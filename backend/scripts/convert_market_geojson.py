"""Convert market GeoJSON using existing layouts as identity/metadata templates.

Run from any directory: python3 backend/scripts/convert_market_geojson.py
Uses the current JSON when present, otherwise the committed JSON template.
Does not modify the database. Original GeoJSON remains the source of truth.
"""
import collections
import argparse
import copy
import json
import math
from pathlib import Path
import re
import subprocess

ROOT = Path(__file__).resolve().parents[2]
DATA = ROOT / 'backend/data/markets'


def ring_key(ring):
    return tuple(sorted(set(tuple(point[:2]) for point in ring)))


def center(ring):
    points = ring[:-1] if ring[0] == ring[-1] else ring
    return [sum(p[i] for p in points) / len(points) for i in (0, 1)]


def distance(a, b):
    return (a[0] - b[0]) ** 2 + (a[1] - b[1]) ** 2


def valid_polygon(feature):
    """Reject closed lines/points that were exported as Polygon geometry."""
    ring = feature['geometry']['coordinates'][0]
    unique_points = set(tuple(point[:2]) for point in ring)
    if len(unique_points) < 3:
        return False
    origin_x, origin_y = ring[0][:2]
    twice_area = abs(sum(
        (point_a[0] - origin_x) * (point_b[1] - origin_y)
        - (point_b[0] - origin_x) * (point_a[1] - origin_y)
        for point_a, point_b in zip(ring, ring[1:])
    ))
    return twice_area >= 1e-16


def hull(points):
    points = sorted(set(tuple(p[:2]) for p in points))
    def cross(o, a, b):
        return (a[0]-o[0])*(b[1]-o[1])-(a[1]-o[1])*(b[0]-o[0])
    def half(seq):
        out = []
        for p in seq:
            while len(out) >= 2 and cross(out[-2], out[-1], p) <= 0:
                out.pop()
            out.append(p)
        return out
    result = half(points)[:-1] + half(reversed(points))[:-1]
    return [list(p) for p in result + result[:1]]


def group(feature):
    p = feature['properties']
    return (str(p.get('name') or ''), str(p.get('zone') or ''), str(p.get('source_layer') or ''))


def number(feature):
    p = feature['properties']
    value = str(p.get('number', p.get('Number')) or '').strip()
    return str(int(value)) if value.isdigit() else value


def stall_number(stall):
    code = stall['code'].split('-')
    value = code[2] if len(code) > 2 else ''
    return str(int(value)) if value.isdigit() else value


def min_buri_facility_type(name):
    if name in ('ห้องน้ำ', 'อ่างล้างมือ'):
        return 'sanitation'
    if name in ('ประตู', 'บันได', 'บันไดเลื่อน'):
        return 'access'
    if name in ('เต้น', 'ลานเร่ เช้า-เย็น', 'ศาลเจ้า', 'ศาลพระพรหม', 'หอนาฬิกา'):
        return 'common-area'
    if name in ('โรงรับจำนำ', 'อาคารกึ่งถาวร'):
        return 'building'
    return 'service'


def min_buri_facility_project(props):
    """Keep the 18 surveyed facilities that belong to a project explicit.

    The remaining facility polygons have zone "-" and are market-level common
    facilities. These keys come from their original project membership before
    facilities were separated from rentable stalls.
    """
    key = (
        str(props.get('source_layer') or ''),
        str(props.get('name') or '').strip(),
        number({'properties': props})
    )
    project_01 = {
        ('2', 'บันไดเลื่อน', '1'), ('2', 'บันไดเลื่อน', '2'),
        ('2', 'โรงรับจำนำ', '1'), ('2', 'บก.', '1'),
        ('2', 'อาคารกึ่งถาวร', '1'), ('2', 'ศาลพระพรหม', '1'),
        ('2', 'อ่างล้างมือ', '1'), ('2', 'อ่างล้างมือ', '2'), ('2', 'อ่างล้างมือ', '3'),
        ('2', 'ลานเร่ เช้า-เย็น', '1'), ('2', 'ลานเร่ เช้า-เย็น', '2')
    }
    project_12 = {
        ('1', 'DIY', '1'), ('1', 'ห้องประชุม', '1'), ('1', 'งานสถานที่', '1'),
        ('1', 'อ่างล้างมือ', '1'), ('1', 'อ่างล้างมือ', '2'),
        ('1', 'อ่างล้างมือ', '3'), ('1', 'อ่างล้างมือ', '4')
    }
    if key in project_01:
        return 'min-01'
    if key in project_12:
        return 'min-12'
    return None


def split_min_buri_facilities(doc, features):
    zone_projects = {
        '1': 'min-06', 'โครงการที่ 1': 'min-06',
        'โครงการที่ 2': 'min-04', 'โครงการที่ 3': 'min-03',
        'โครงการที่ 4': 'min-05', 'โครงการที่ 5': 'min-10',
        'โครงการที่ 6': 'min-08', 'โครงการที่ 7': 'min-07',
        'โครงการที่ 8': 'min-09', 'แผงค้าเบ็ดเตล็ด': 'min-02',
        'แผงเบ็ดเตล็ด': 'min-02', 'แผงค้าต้นไม้': 'min-11'
    }
    facilities = []
    facility_index = 0
    for feature in features:
        props = feature['properties']
        name = str(props.get('name') or '').strip()
        if name == 'แผงร้านค้า':
            continue
        facility_index += 1
        source_layer = str(props.get('source_layer') or '')
        floor_id = '2' if source_layer == '1' else '1'
        source_zone = str(props.get('zone') or '').strip()
        outer_ring = feature['geometry']['coordinates'][0]
        facility = {
            'id': f'min-buri-facility-{facility_index:03d}',
            'name': name,
            'type': min_buri_facility_type(name),
            'coordinates': center(outer_ring),
            'geometry': feature['geometry'],
            'floorId': floor_id,
            'sourceZone': source_zone,
            'number': number(feature),
            'sourceProperties': props
        }
        project_id = zone_projects.get(source_zone) or min_buri_facility_project(props)
        if project_id:
            facility['projectId'] = project_id
        facilities.append(facility)

    for project in doc['projects']:
        project['stalls'] = [
            stall for stall in project.get('stalls', [])
            if str(stall.get('sourceProperties', {}).get('name') or '').strip() == 'แผงร้านค้า'
        ]
        project['totalStalls'] = len(project['stalls'])
        project['stallsAvailable'] = sum(stall.get('status') == 'available' for stall in project['stalls'])
        project['realAreaSqm'] = round(sum(stall.get('size', 0) for stall in project['stalls']), 2)
        if project['stalls']:
            project['coordinates'] = hull([
                point for stall in project['stalls'] for point in stall['coordinates']
            ])
            project['centroid'] = center(project['coordinates'])
    doc['facilities'] = facilities
    doc['siteFeatures'] = []
    return facilities


def normalize_min_buri_project_names(doc):
    """Use the surveyed zone as the displayed project name.

    IDs such as min-05 remain stable internal references. They must not be
    presented as the real project number when the source says โครงการที่ 4.
    """
    doc['projects'] = [project for project in doc['projects'] if project.get('stalls')]
    for project in doc['projects']:
        source_zones = collections.Counter(
            str(stall.get('sourceProperties', {}).get('zone') or '').strip()
            for stall in project.get('stalls', [])
            if str(stall.get('sourceProperties', {}).get('zone') or '').strip()
        )
        floor_name = (project.get('floors') or [{}])[0].get('name') or 'ชั้นล่าง'
        source_zone = source_zones.most_common(1)[0][0] if source_zones else 'โซนหลัก'
        project['name'] = f'{source_zone} ({floor_name})'

        source_project_number = re.fullmatch(r'โครงการที่\s*(\d+)', source_zone)
        if source_project_number:
            project['sectionNumber'] = source_project_number.group(1)
        else:
            # Do not expose generated IDs (01, 02, 11) as real project
            # numbers when the survey contains only a zone name or no zone.
            project['sectionNumber'] = ''

        for stall in project.get('stalls', []):
            stall['zone'] = project['name']


def convert(path, committed_template=False):
    output = path.with_suffix('.json')
    raw = output.read_text() if output.exists() and not committed_template else subprocess.check_output(
        ['git', 'show', 'HEAD:' + str(output.relative_to(ROOT))], cwd=ROOT, text=True)
    doc = json.loads(raw)
    source_features = json.loads(path.read_text())['features']
    for f in source_features:
        g = f['geometry']
        assert g['type'] == 'Polygon' and g['coordinates'], 'Unsupported polygon'
        for ring in g['coordinates']:
            assert len(ring) >= 4 and ring[0] == ring[-1], 'Invalid ring'
            assert all(len(p) >= 2 and all(math.isfinite(v) for v in p[:2]) for p in ring)
    invalid_features = [feature for feature in source_features if not valid_polygon(feature)] if path.stem == 'min-buri' else []
    features = [feature for feature in source_features if valid_polygon(feature)] if path.stem == 'min-buri' else source_features
    projects = doc['projects']
    old = [(p, s) for p in projects for s in p.get('stalls', [])]
    lookup = collections.defaultdict(list)
    for index, (_, s) in enumerate(old):
        lookup[ring_key(s['coordinates'])].append(index)
    matches, used, votes = {}, set(), collections.defaultdict(collections.Counter)
    for i, f in enumerate(features):
        candidates = [j for j in lookup[ring_key(f['geometry']['coordinates'][0])] if j not in used]
        if candidates:
            j = candidates[0]
            matches[i] = j
            used.add(j)
            votes[group(f)][old[j][0]['id']] += 1
    group_projects = {key: counts.most_common(1)[0][0] for key, counts in votes.items()}
    if path.stem == 'bang-kapi':
        # Bang Kapi's survey schema is explicit: `name` is the project name
        # and `number` is the stall name. Do not classify office/toilet/waste
        # polygons as facilities or infer their project from geometry.
        group_projects = {}
        for key in dict.fromkeys(group(f) for f in features):
            source_name = key[0]
            candidates = [p for p in projects if source_name in p['name']]
            if len(candidates) == 1:
                group_projects[key] = candidates[0]['id']
    for f in features:
        key = group(f)
        if key in group_projects:
            continue
        name, zone, layer = key
        rental = path.stem == 'bang-kapi' or any(word in name for word in ('แผง', 'ร้านค้า', 'อาคารพาณิชย์', 'พื้นที่ด้านซอย', 'ร้านสหกรณ์'))
        if not rental:
            continue
        label = zone if zone not in ('', '-') else (name if path.stem == 'bang-kapi' else 'โซนหลัก')
        candidates = [p for p in projects if '(' + label + ')' in p['name'] or '(' + label + ' (' in p['name']]
        if layer and path.stem == 'min-buri':
            floor = '2' if layer == '1' else '1'
            candidates = [p for p in candidates if str(p['floors'][0]['id']) == floor]
        if len(candidates) == 1:
            group_projects[key] = candidates[0]['id']
    by_id = {p['id']: p for p in projects}
    new_projects = []
    assignments = {}
    for i, f in enumerate(features):
        key = group(f)
        if path.stem == 'bang-kapi' and key in group_projects:
            assignments[i] = group_projects[key]
        elif i in matches:
            assignments[i] = old[matches[i]][0]['id']
        elif key in group_projects:
            assignments[i] = group_projects[key]
        elif path.stem == 'bang-kapi' or any(word in key[0] for word in ('แผง', 'ร้านค้า', 'อาคารพาณิชย์')):
            next_number = len(projects) + len(new_projects) + 1
            project_id = f"ban-{next_number:02d}" if path.stem == 'bang-kapi' else 'geo-' + str(len(new_projects) + 1)
            while project_id in by_id:
                project_id += '-new'
            label = key[1] if key[1] not in ('', '-') else key[0]
            p = copy.deepcopy(projects[0])
            section_number = f"{next_number:02d}" if path.stem == 'bang-kapi' else 'G' + str(len(new_projects)+1)
            p.update(id=project_id, name=label, sectionNumber=section_number, category=label, stalls=[])
            by_id[project_id] = p
            new_projects.append(p)
            group_projects[key] = project_id
            assignments[i] = project_id
    # Match remaining records within each project, preferring the source number.
    pairs = []
    for i, pid in assignments.items():
        if i in matches:
            continue
        f = features[i]
        c = center(f['geometry']['coordinates'][0])
        for j, (p, s) in enumerate(old):
            if j in used or p['id'] != pid:
                continue
            same_number = bool(number(f)) and number(f) == stall_number(s)
            d = distance(c, center(s['coordinates']))
            # Unnumbered records must remain within 30 m of their original location.
            if same_number or d < (30 / 111000) ** 2:
                pairs.append((not same_number, d, i, j))
    for _, _, i, j in sorted(pairs):
        if i not in matches and j not in used:
            matches[i] = j
            used.add(j)
    templates = {p['id']: copy.deepcopy(p['stalls'][0]) for p in projects if p.get('stalls')}
    all_projects = projects + new_projects
    for p in all_projects:
        p['stalls'] = []
    for i, pid in assignments.items():
        f = features[i]
        p = by_id[pid]
        if i in matches and old[matches[i]][0]['id'] == pid:
            s = copy.deepcopy(old[matches[i]][1])
        else:
            s = copy.deepcopy(templates.get(pid, next(iter(templates.values()))))
            if path.stem == 'bang-kapi':
                source_number = number(f)
                code_suffix = source_number.zfill(3) if source_number.isdigit() else source_number
                code = f"BKA-{p['sectionNumber']}-{code_suffix}"
            else:
                code = f"{doc['slug']}-GEO-{i+1:05}"
            s.update(code=code, status='available', projectId=pid,
                     zone=p['name'], floor=str(p['floors'][0]['id']),
                     floorName=p['floors'][0]['name'], floorLevel=p['floors'][0]['level'])
        if path.stem == 'bang-kapi':
            s['name'] = number(f)
            s['projectId'] = pid
            s['zone'] = p['name']
        s['coordinates'] = f['geometry']['coordinates'][0]
        s['polygonHoles'] = f['geometry']['coordinates'][1:]
        s['center'] = center(s['coordinates'])
        s['sourceProperties'] = f['properties']
        p['stalls'].append(s)
    for p in all_projects:
        if path.stem == 'bang-kapi' and p['stalls']:
            source_name = p['stalls'][0]['sourceProperties']['name']
            p['name'] = source_name
            p['category'] = source_name
            for stall in p['stalls']:
                stall['zone'] = source_name
        p['totalStalls'] = len(p['stalls'])
        p['stallsAvailable'] = sum(s.get('status') == 'available' for s in p['stalls'])
        p['realAreaSqm'] = round(sum(s.get('size', 0) for s in p['stalls']), 2)
        if p['stalls']:
            p['coordinates'] = hull([pt for s in p['stalls'] for pt in s['coordinates']])
            p['centroid'] = center(p['coordinates'])
    doc['projects'] = all_projects
    # Keep non-rental polygons intact without turning toilets/offices into rentable stalls.
    doc['siteFeatures'] = [f for i, f in enumerate(features) if i not in assignments]
    facilities = split_min_buri_facilities(doc, features) if path.stem == 'min-buri' else []
    if path.stem == 'min-buri':
        normalize_min_buri_project_names(doc)
        doc['invalidFeatures'] = invalid_features
    codes = [s['code'] for p in all_projects for s in p['stalls']]
    assert len(codes) == len(set(codes)), 'Duplicate stall identifiers'
    if path.stem == 'min-buri':
        assert len(codes) + len(facilities) + len(invalid_features) == len(source_features)
    else:
        assert len(assignments) + len(doc['siteFeatures']) == len(features)
    return output, doc, {'market': path.stem, 'features': len(source_features), 'stalls': len(codes),
                         'retained_codes': len(matches), 'new_codes': len(assignments)-len(matches),
                         'unmatched_old': len(old)-len(used), 'facilities': len(facilities),
                         'site_features': len(doc['siteFeatures']), 'invalid_features': len(invalid_features)}


def convert_source_schema(path, committed_template=False):
    """Map every source polygon using `name` as project and `number` as stall name."""
    output = path.with_suffix('.json')
    raw = output.read_text() if output.exists() and not committed_template else subprocess.check_output(
        ['git', 'show', 'HEAD:' + str(output.relative_to(ROOT))], cwd=ROOT, text=True)
    doc = json.loads(raw)
    features = json.loads(path.read_text())['features']
    old_projects = doc['projects']
    old_stalls = [(project, stall) for project in old_projects for stall in project.get('stalls', [])]

    geometry_lookup = collections.defaultdict(list)
    for old_index, (_, stall) in enumerate(old_stalls):
        geometry_lookup[ring_key(stall['coordinates'])].append(old_index)

    matches = {}
    votes = collections.defaultdict(collections.Counter)
    used_old = set()
    for feature_index, feature in enumerate(features):
        geometry = feature['geometry']
        assert geometry['type'] == 'Polygon' and geometry['coordinates'], 'Unsupported polygon'
        for ring in geometry['coordinates']:
            assert len(ring) >= 4 and ring[0] == ring[-1], 'Invalid ring'
            assert all(len(point) >= 2 and all(math.isfinite(value) for value in point[:2]) for point in ring)
        candidates = [index for index in geometry_lookup[ring_key(geometry['coordinates'][0])] if index not in used_old]
        if candidates:
            source_props = feature['properties']
            exact = [index for index in candidates if old_stalls[index][1].get('sourceProperties') == source_props]
            old_index = (exact or candidates)[0]
            matches[feature_index] = old_index
            used_old.add(old_index)
            source_name = str(source_props.get('name') or '').strip()
            votes[source_name][old_stalls[old_index][0]['id']] += 1

    source_names = list(dict.fromkeys(str(feature['properties'].get('name') or '').strip() for feature in features))
    templates_by_id = {project['id']: project for project in old_projects}
    used_project_ids = set()
    projects_by_name = {}
    palette = ['#2563eb', '#7c3aed', '#16a34a', '#d97706', '#dc2626', '#0891b2', '#9333ea', '#475569']

    for position, source_name in enumerate(source_names, 1):
        selected = None
        for project_id, _ in votes[source_name].most_common():
            if project_id not in used_project_ids:
                selected = copy.deepcopy(templates_by_id[project_id])
                break
        if selected is None:
            selected = copy.deepcopy(old_projects[0])
            project_id = f"source-{position:02d}"
            while project_id in used_project_ids or project_id in templates_by_id:
                project_id += '-new'
            selected['id'] = project_id
            selected['color'] = palette[(position - 1) % len(palette)]
        used_project_ids.add(selected['id'])
        display_name = source_name or 'ไม่ระบุชื่อโครงการ'
        selected.update(
            name=display_name,
            sourceName=source_name,
            sourceMapped=True,
            sectionNumber=f"{position:02d}",
            category=display_name,
            stalls=[]
        )
        projects_by_name[source_name] = selected

    global_stall_template = copy.deepcopy(old_stalls[0][1])
    retained_codes = 0
    for feature_index, feature in enumerate(features):
        props = feature['properties']
        source_name = str(props.get('name') or '').strip()
        project = projects_by_name[source_name]
        old_match = matches.get(feature_index)
        if old_match is not None:
            stall = copy.deepcopy(old_stalls[old_match][1])
            retained_codes += 1
        else:
            stall = copy.deepcopy(global_stall_template)
            stall.update(
                code=f"{doc['slug'].upper()}-GEO-{feature_index + 1:05d}",
                status='available'
            )
            layer = str(props.get('source_layer') or '')
            if path.stem == 'min-buri' and layer in ('1', '2'):
                floor_id = '2' if layer == '1' else '1'
                stall.update(
                    floor=floor_id,
                    floorName='ชั้นบน' if floor_id == '2' else 'ชั้นล่าง',
                    floorLevel=int(floor_id)
                )

        source_zone = str(props.get('zone') or '').strip()
        stall.update(
            name=number(feature),
            projectId=project['id'],
            zone=source_zone if source_zone not in ('', '-') else project['name'],
            sourceZone=source_zone,
            coordinates=feature['geometry']['coordinates'][0],
            polygonHoles=feature['geometry']['coordinates'][1:],
            center=center(feature['geometry']['coordinates'][0]),
            sourceProperties=props
        )
        project['stalls'].append(stall)

    all_projects = list(projects_by_name.values())
    for project in all_projects:
        floor_metadata = {}
        template = templates_by_id.get(project['id'], {})
        for floor in template.get('floors', []):
            floor_metadata[str(floor['id'])] = copy.deepcopy(floor)
        for stall in project['stalls']:
            floor_id = str(stall.get('floor') or stall.get('floorLevel') or '1')
            stall['floor'] = floor_id
            if floor_id not in floor_metadata:
                level = int(floor_id) if floor_id.isdigit() else len(floor_metadata) + 1
                floor_metadata[floor_id] = {
                    'id': floor_id,
                    'name': stall.get('floorName') or f'ชั้น {level}',
                    'shortLabel': f'{level}F',
                    'level': level
                }
        project['floors'] = sorted(floor_metadata.values(), key=lambda floor: floor.get('level', 0))
        project['totalStalls'] = len(project['stalls'])
        project['stallsAvailable'] = sum(stall.get('status') == 'available' for stall in project['stalls'])
        project['realAreaSqm'] = round(sum(stall.get('size', 0) for stall in project['stalls']), 2)
        project['coordinates'] = hull([
            point for stall in project['stalls'] for point in stall['coordinates']
        ])
        project['centroid'] = center(project['coordinates'])

    doc['projects'] = all_projects
    doc['siteFeatures'] = []
    codes = [stall['code'] for project in all_projects for stall in project['stalls']]
    assert len(codes) == len(set(codes)), 'Duplicate internal stall identifiers'
    assert len(codes) == len(features), 'Every source polygon must map to a stall record'
    return output, doc, {
        'market': path.stem,
        'features': len(features),
        'projects': len(all_projects),
        'retained_codes': retained_codes,
        'new_codes': len(features) - retained_codes,
        'site_features': 0
    }


if __name__ == '__main__':
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument('--committed-template', action='store_true', help='Use HEAD layouts as templates')
    parser.add_argument('--market', help='Convert only this market slug')
    args = parser.parse_args()
    paths = [DATA / f'{args.market}.geojson'] if args.market else sorted(DATA.glob('*.geojson'))
    converted = [convert(path, args.committed_template) for path in paths]
    for path, doc, summary in converted:
        path.write_text(json.dumps(doc, ensure_ascii=False, indent=2) + '\n')
        print(json.dumps(summary, ensure_ascii=False))
