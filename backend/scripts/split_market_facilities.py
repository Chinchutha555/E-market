#!/usr/bin/env python3
"""Split non-rental polygons into `facilities`, the same shape Min Buri uses.

Min Buri keeps rentable stalls in `projects[].stalls` and everything else
(toilets, gates, offices, buildings...) in the top-level `facilities` list.
This applies the same split to the other markets without regenerating stall
codes: it moves the matching stalls and every `siteFeatures` entry into
`facilities`, then recomputes the project totals.

Usage: python3 scripts/split_market_facilities.py [--dry-run] [slug ...]
Idempotent: markets that already have `facilities` are skipped.
"""
import json
import math
import sys
from pathlib import Path

MARKETS = Path(__file__).resolve().parent.parent / 'data' / 'markets'

# Names (found on stalls) that are not rentable stalls. siteFeatures are always facilities.
NON_RENTAL_STALL_NAMES = {
    'bang-kapi': {'อาคารพาณิชย์', 'พื้นที่ด้านซอย 125', 'พื้นที่ด้านซอย 127', 'ห้องสุขา', 'ที่พักขยะ', 'ที่ทำการ'},
    'chatuchak': {'ATM'},
    'nong-chok': {'ร้านสหกรณ์', 'ศาลาภู่อารีย์', 'สะพานข้ามคลองแสนแสบ', 'อาคารพาณิชย์'},
    'pracha-niwet-1': {'จุดบริการไฟ', 'อาคารโกล์ดมาร์เก็ต', 'พื้นที่ว่าง', 'โยธา'},
    'ratchadaphisek': {'ซักล้าง', 'เก็บของ', 'สนามกีฬา', 'ท่าเรือ'},
    'thewarat': {'อาคารพาณิชย์'},
    'thonburi-sanamluang-2': {'โรงน้ำแข็ง', 'JEN', 'คลังโยธา', 'อาคารกองอำนวยการ', 'บ้านพัก'},
}

SANITATION = ('ห้องน้ำ', 'ห้องสุขา', 'สุขา', 'อ่างล้างมือ', 'ซักล้าง', 'ที่พักขยะ', 'ระบบบำบัด', 'บ่อบำบัด', 'แท้งค์น้ำ')
ACCESS = ('ประตู', 'บันได', 'สะพาน', 'ท่าน้ำ', 'ท่าเรือ', 'ที่จอดรถ', 'ลานจอดรถ', 'พื้นที่ด้านซอย')
COMMON = ('ลาน', 'ศาล', 'พระพรหม', 'หอนาฬิกา', 'สนามกีฬา', 'สวนสาธารณะ', 'พื้นที่ว่าง', 'เต้น')
BUILDING = ('อาคาร', 'ตึกแถว', 'บ้านพัก', 'โรงน้ำแข็ง', 'โรงรับจำนำ', 'สถานธนานุบาล', 'JJ Mall', 'จตุจักรพลาซ่า', 'ร้านอาหาร', 'JEN', 'ร้านสหกรณ์')


def facility_type(name):
    if name.startswith('ลานจอดรถ') or name.startswith('ที่จอดรถ'):
        return 'access'
    for prefixes, kind in ((SANITATION, 'sanitation'), (ACCESS, 'access'), (COMMON, 'common-area'), (BUILDING, 'building')):
        if name.startswith(prefixes):
            return kind
    return 'service'


def ring_center(ring):
    points = ring[:-1] if ring[0] == ring[-1] else ring
    return [sum(p[i] for p in points) / len(points) for i in (0, 1)]


def hull(points):
    pts = sorted(set(map(tuple, points)))
    if len(pts) <= 2:
        return [list(p) for p in pts]

    def cross(o, a, b):
        return (a[0] - o[0]) * (b[1] - o[1]) - (a[1] - o[1]) * (b[0] - o[0])

    def half(seq):
        out = []
        for p in seq:
            while len(out) >= 2 and cross(out[-2], out[-1], p) <= 0:
                out.pop()
            out.append(p)
        return out

    lower, upper = half(pts), half(reversed(pts))
    return [list(p) for p in lower[:-1] + upper[:-1]]


def facility_from(doc, index, name, geometry, floor_id, source_zone, number, props, project_id=None):
    outer = geometry['coordinates'][0]
    facility = {
        'id': f"{doc['slug']}-facility-{index:03d}",
        'name': name,
        'type': facility_type(name),
        'coordinates': ring_center(outer),
        'geometry': geometry,
        'floorId': floor_id,
        'sourceZone': source_zone,
        'number': number,
        'sourceProperties': props,
    }
    if project_id:
        facility['projectId'] = project_id
    return facility


def split(doc):
    slug = doc['slug']
    non_rental = NON_RENTAL_STALL_NAMES.get(slug, set())
    facilities = []
    moved_codes = []

    for project in doc['projects']:
        kept = []
        for stall in project.get('stalls', []):
            props = stall.get('sourceProperties') or {}
            name = str(props.get('name') or '').strip()
            if name not in non_rental:
                kept.append(stall)
                continue
            geometry = {'type': 'Polygon', 'coordinates': [stall['coordinates']] + (stall.get('polygonHoles') or [])}
            facilities.append(facility_from(
                doc, len(facilities) + 1, name, geometry, str(stall.get('floor') or '1'),
                str(props.get('zone') or '').strip() or '-', props.get('number', stall.get('name')),
                props, project.get('id')))
            moved_codes.append(stall['code'])
        project['stalls'] = kept
        project['totalStalls'] = len(kept)
        project['stallsAvailable'] = sum(s.get('status') == 'available' for s in kept)
        project['realAreaSqm'] = round(sum(s.get('size', 0) for s in kept), 2)
        if kept:
            project['coordinates'] = hull([pt for s in kept for pt in s['coordinates']])
            project['centroid'] = ring_center(project['coordinates'])

    for feature in doc.get('siteFeatures', []):
        props = feature.get('properties') or {}
        name = str(props.get('name') or '').strip()
        number = props.get('number', props.get('Number'))
        facilities.append(facility_from(
            doc, len(facilities) + 1, name, feature['geometry'], '1',
            str(props.get('zone') or '').strip() or '-', number, props))

    # A "project" that held only non-rental polygons is just a facility: drop it so the
    # facility is not covered by an empty project polygon that swallows map clicks.
    dropped = {p['id'] for p in doc['projects'] if not p.get('stalls')}
    if dropped:
        doc['projects'] = [p for p in doc['projects'] if p['id'] not in dropped]
        for facility in facilities:
            if facility.get('projectId') in dropped:
                del facility['projectId']

    doc['facilities'] = facilities
    doc['siteFeatures'] = []
    return facilities, moved_codes


def main():
    args = [a for a in sys.argv[1:] if not a.startswith('--')]
    dry = '--dry-run' in sys.argv
    for path in sorted(MARKETS.glob('*.json')):
        if args and path.stem not in args:
            continue
        doc = json.loads(path.read_text())
        if 'facilities' in doc:
            print(f'{path.stem}: skip (already has facilities)')
            continue
        before = sum(len(p.get('stalls', [])) for p in doc['projects'])
        facilities, moved = split(doc)
        after = sum(len(p.get('stalls', [])) for p in doc['projects'])
        assert before - after == len(moved)
        codes = [s['code'] for p in doc['projects'] for s in p['stalls']]
        assert len(codes) == len(set(codes))
        kinds = {}
        for f in facilities:
            kinds[f['type']] = kinds.get(f['type'], 0) + 1
        print(f"{path.stem}: stalls {before} -> {after}, facilities {len(facilities)} {kinds}")
        if not dry:
            path.write_text(json.dumps(doc, ensure_ascii=False, indent=2) + '\n')


if __name__ == '__main__':
    main()
