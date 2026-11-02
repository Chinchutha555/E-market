<template>
  <div class="projects-master-container" :class="{ 'is-preview-mode': previewMode, 'is-management': managementMode }">
    <!-- Top Interactive Filter & Quick Search Bar -->
    <Teleport :to="canTeleportToolbar ? '#step1-map-toolbar' : 'body'" :disabled="!canTeleportToolbar">
    <div v-if="!previewMode && !isOverviewPage" class="master-map-top-bar">
      <div class="master-map-tools-right">
        <!-- Search Input with Clear Button -->
        <div class="search-input-group">
          <span class="search-icon">
            <svg width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
              <circle cx="11" cy="11" r="8"></circle>
              <line x1="21" y1="21" x2="16.65" y2="16.65"></line>
            </svg>
          </span>
          <input
            v-model="searchQuery"
            type="text"
            :placeholder="managementMode ? 'เลขที่การจอง คำร้อง ชื่อโครงการ หรือรหัสแผง' : 'ชื่อโครงการหรือประเภทสินค้า'"
            class="map-search-input"
            @input="handleSearchInput"
            @keydown.down.prevent="moveSearchSuggestion(1)"
            @keydown.up.prevent="moveSearchSuggestion(-1)"
            @keydown.enter.prevent="selectActiveSearchSuggestion"
          />
          <button
            v-if="searchQuery"
            type="button"
            class="clear-search-btn"
            @click="clearSearch"
            aria-label="ล้างคำค้นหา"
          >
            <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><line x1="18" y1="6" x2="6" y2="18"></line><line x1="6" y1="6" x2="18" y2="18"></line></svg>
          </button>
        </div>

        <div class="map-advanced-filter-wrap">
          <button type="button" class="map-advanced-filter-btn" :class="{ active: hasAdvancedFilters }" @click="showAdvancedFilters = !showAdvancedFilters" :aria-expanded="showAdvancedFilters">
            <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M4 6h16M7 12h10m-7 6h4"></path></svg>
            ตัวกรอง
            <span v-if="advancedFilterCount" class="map-advanced-count">{{ advancedFilterCount }}</span>
          </button>
          <div v-if="showAdvancedFilters" class="map-advanced-filter-panel">
            <div class="map-filter-panel-head">
              <strong>ตัวกรองเพิ่มเติม</strong>
              <button type="button" @click="showAdvancedFilters = false" aria-label="ปิดตัวกรอง">
                <svg width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                  <line x1="18" y1="6" x2="6" y2="18"></line>
                  <line x1="6" y1="6" x2="18" y2="18"></line>
                </svg>
              </button>
            </div>

            <div v-if="managementMode" class="filter-field-group">
              <span class="filter-field-label">โซน / โครงการ</span>
              <AppSelect
                v-model="zoneFilter"
                :options="zoneOptions"
                placeholder="ทุกโซน"
              />
            </div>

            <div class="filter-field-group">
              <span class="filter-field-label">ประเภทสินค้า</span>
              <AppSelect
                :model-value="draftCategoryFilter"
                :options="categoryOptions"
                placeholder="เลือกหมวดหมู่"
                @update:model-value="handleCategoryChange"
              />
            </div>

            <div class="filter-field-group">
              <span class="filter-field-label">งบค่าเช่าต่อเดือน</span>
              <AppSelect
                v-model="rentFilter"
                :options="rentOptions"
                placeholder="งบค่าเช่าต่อเดือน"
              />
            </div>

            <div class="filter-field-group">
              <span class="filter-field-label">สถานะแผง</span>
              <AppSelect
                v-model="availabilityFilter"
                :options="availabilityOptions"
                placeholder="สถานะแผง"
              />
            </div>

            <div v-if="managementMode" class="filter-field-group">
              <span class="filter-field-label">ความพร้อมของข้อมูล</span>
              <AppSelect
                v-model="readinessFilter"
                :options="readinessOptions"
                placeholder="ทั้งหมด"
              />
            </div>

            <div class="map-filter-panel-footer">
              <div class="mfp-result-summary">
                <span class="mfp-result-label">ผลลัพธ์:</span>
                <span class="mfp-result-val"><strong>{{ previewFilteredProjectCount }}</strong> โครงการ</span>
                <span class="mfp-result-dot">·</span>
                <span class="mfp-result-stalls"><strong>{{ previewFilteredStallCount.toLocaleString() }}</strong> แผง</span>
              </div>
              <div class="mfp-action-grid">
                <button
                  type="button"
                  class="btn-filter-reset"
                  @click="clearAllFilters"
                  title="รีเซ็ตตัวกรองทั้งหมด"
                >
                  <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.2" stroke-linecap="round" stroke-linejoin="round">
                    <path d="M3 12a9 9 0 1 0 9-9 9.75 9.75 0 0 0-6.74 2.74L3 8"></path>
                    <path d="M3 3v5h5"></path>
                  </svg>
                  <span>ล้างค่า</span>
                </button>
                <button
                  type="button"
                  class="btn-filter-apply"
                  @click="applyFilters"
                  title="ดูโครงการที่ตรงตามเงื่อนไข"
                >
                  <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.2" stroke-linecap="round" stroke-linejoin="round">
                    <circle cx="11" cy="11" r="8"></circle>
                    <line x1="21" y1="21" x2="16.65" y2="16.65"></line>
                  </svg>
                  <span>ค้นหา ({{ previewFilteredProjectCount }})</span>
                </button>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
    </Teleport>

    <!-- Preview Mode Top Bar (สำหรับขั้นตอนตรวจสอบยืนยันข้อมูล) -->
    <div v-if="previewMode" class="preview-top-bar">
      <!-- ด้านซ้าย: หัวข้อแผนที่ + ไอคอน + รายละเอียดโครงการที่จอง -->
      <div class="preview-head-left">
        <div class="preview-head-icon-box">
          <svg width="19" height="19" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.2" stroke-linecap="round" stroke-linejoin="round">
            <polygon points="1 6 1 22 8 18 16 22 23 18 23 2 16 6 8 2 1 6"></polygon>
            <line x1="8" y1="2" x2="8" y2="18"></line>
            <line x1="16" y1="6" x2="16" y2="22"></line>
          </svg>
        </div>
        <div class="preview-head-text-col">
          <div class="preview-head-title-row">
            <strong class="preview-head-title">ตำแหน่งแผงค้าบนแผนที่จริง</strong>
            <span v-if="bookedProject" class="preview-head-cat-badge">
              {{ bookedProject.name }}
            </span>
          </div>
          <div class="preview-head-subtitle">
            {{ bookedProject?.category || 'ตรวจสอบตำแหน่งแผงค้าที่คุณเลือกจองในผังตลาด' }}
          </div>
        </div>
      </div>

      <!-- ด้านขวา: ปุ่มโฟกัสแผงที่จอง + แถบอธิบายสีแผง (Legend) -->
      <div class="preview-head-right">
        <!-- Clickable Booked Stalls Button / Pill -->
        <button
          type="button"
          class="preview-booked-pill-btn"
          :class="{ 'is-active': activeProjectId === bookedProjectId }"
          @click="focusBookedStalls"
          title="คลิกเพื่อซูมกลับไปยังแผงที่คุณเลือกจองทันที"
        >
          <span class="preview-info-dot" :style="{ backgroundColor: bookedProject?.color || '#2563eb' }"></span>
          <span class="preview-badge-stalls">
            แผงที่คุณเลือก {{ selectedStallCodes.length }} แผง ({{ selectedStallCodes.join(', ') }})
          </span>
          <span v-if="isNotViewingBookedProject" class="return-hint-tag">
            <svg width="11" height="11" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round" style="margin-right: 3px;">
              <polyline points="1 4 1 10 7 10"></polyline>
              <path d="M3.51 15a9 9 0 1 0 2.13-9.36L1 10"></path>
            </svg>
            คลิกเพื่อกลับมาดู
          </span>
        </button>

        <span v-if="activeProjectId && isNotViewingBookedProject" class="viewing-other-badge" style="font-size: 0.78rem; color: var(--tx-475569); background: var(--bg-f1f5f9); padding: 4px 10px; border-radius: 6px; display: inline-flex; align-items: center; gap: 5px;">
          <span>กำลังดู:</span>
          <strong style="color: var(--tx-0f172a);">{{ activeProject?.name }}</strong>
        </span>

        <div class="preview-legend-row">
          <div class="preview-leg-item"><span class="leg-dot is-selected"></span><span>แผงที่คุณเลือก ({{ selectedStallCodes.length }})</span></div>
          <div class="preview-leg-item"><span class="leg-dot is-available"></span><span>แผงว่าง</span></div>
          <div class="preview-leg-item"><span class="leg-dot is-reserved"></span><span>ติดจอง</span></div>
        </div>
      </div>
    </div>

    <!-- MapLibre GL JS Interactive Map Container -->
    <div class="maplibre-overview-wrapper" :class="{ 'is-drawing-active': isDrawingZone || isDrawingStall || isAddingFacility, 'is-map-loading': isMapLoading, 'has-ov-panel': isOverviewPage && !!ovPanel && ovPanel !== 'dashboard', 'has-ov-dashboard': isOverviewPage && ovPanel === 'dashboard' }">
      <div :id="containerId" ref="mapContainerRef" class="maplibre-overview-canvas"></div>

      <div v-if="isMapLoading" class="map-loading-skeleton" aria-busy="true" aria-label="กำลังเตรียมผังตลาด">
        <div class="map-skeleton-control-group" aria-hidden="true">
          <span class="map-skeleton-line is-map-tab"></span>
          <span class="map-skeleton-line is-map-tab short"></span>
          <span class="map-skeleton-line is-center-control"></span>
        </div>
        <div class="map-skeleton-zone zone-one" aria-hidden="true"></div>
        <div class="map-skeleton-zone zone-two" aria-hidden="true"></div>
        <div class="map-skeleton-zone zone-three" aria-hidden="true"></div>
        <div class="map-skeleton-legend" aria-hidden="true">
          <span class="map-skeleton-line is-legend-title"></span>
          <span class="map-skeleton-line is-legend-row"></span>
          <span class="map-skeleton-line is-legend-row shorter"></span>
        </div>
        <div class="map-skeleton-tools" aria-hidden="true">
          <span v-for="index in 4" :key="index" class="map-skeleton-tool"></span>
        </div>
      </div>

      <div
        v-if="managementMode && !inspectOnly && !activeProject && !isDrawingZone && !isDrawingStall && !isAddingFacility && !isDeletingStall && !isPolygonEditMode && !isStallPolygonEditMode"
        class="map-project-guidance map-staff-guidance"
        role="status"
      >
        <div class="map-project-guidance-icon" aria-hidden="true">
          <svg v-if="inspectOnly" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.2" stroke-linecap="round" stroke-linejoin="round">
            <path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"></path>
            <polyline points="14 2 14 8 20 8"></polyline>
            <line x1="8" y1="13" x2="16" y2="13"></line>
            <line x1="8" y1="17" x2="13" y2="17"></line>
          </svg>
          <svg v-else width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.2" stroke-linecap="round" stroke-linejoin="round">
            <polygon points="1 6 1 22 8 18 16 22 23 18 23 2 16 6 8 2 1 6"></polygon>
            <line x1="8" y1="2" x2="8" y2="18"></line>
            <line x1="16" y1="6" x2="16" y2="22"></line>
          </svg>
        </div>
        <span>
          <strong>{{ inspectOnly ? 'เลือกโซนเพื่อดูข้อมูลแผง' : 'เลือกโซนเพื่อจัดการผัง' }}</strong>
          <small>{{ inspectOnly ? 'คลิกพื้นที่โซนบนแผนที่เพื่อดูข้อมูลแผงและสัญญาเช่า' : 'คลิกพื้นที่โซนเพื่อจัดการแผง หรือใช้เครื่องมือด้านล่าง' }}</small>
        </span>
      </div>

      <div
        v-if="managementMode && !inspectOnly"
        class="zone-management-toolbar"
        :class="{ 'is-drawing-toolbar': isDrawingZone || isDrawingStall || isAddingFacility || isDeletingStall || isPolygonEditMode || isStallPolygonEditMode || drawingPoints.length > 0 }"
        aria-label="จัดการโซนบนผัง"
      >
        <template v-if="!isDrawingZone && !isDrawingStall && !isAddingFacility && !isDeletingStall && !isPolygonEditMode && !isStallPolygonEditMode && drawingPoints.length === 0">
          <AppButton size="md" variant="outline" @click="beginDrawingZone"><template #prefix><svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><polygon points="9 2 15 2 19 6 19 12 15 16 9 16 5 12 5 6"/><path d="M19 19v4M17 21h4"/></svg></template>เพิ่มโซน</AppButton>
          <AppButton size="md" variant="outline" :disabled="!activeProject" @click="beginDrawingStall"><template #prefix><svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><rect x="4" y="4" width="16" height="16" rx="2"/><path d="M12 8v8M8 12h8"/></svg></template>เพิ่มแผง</AppButton>
          <AppButton size="md" variant="outline" :disabled="!activeProject || !(activeLayout?.stalls?.length)" @click="toggleStallPolygonEdit"><template #prefix><svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><rect x="4" y="4" width="16" height="16" rx="2"/><circle cx="4" cy="4" r="2"/><circle cx="20" cy="4" r="2"/><circle cx="20" cy="20" r="2"/><circle cx="4" cy="20" r="2"/></svg></template>ปรับขนาดแผง</AppButton>
          <AppButton size="md" variant="danger-outline" :disabled="!activeProject || !(activeLayout?.stalls?.length)" @click="beginDeletingStall"><template #prefix><svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><rect x="4" y="4" width="16" height="16" rx="2"/><path d="m9 9 6 6M15 9l-6 6"/></svg></template>ลบแผง</AppButton>
          <AppButton size="md" variant="outline" @click="beginAddingFacility"><template #prefix><svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M12 22s7-5.3 7-12a7 7 0 1 0-14 0c0 6.7 7 12 7 12Z"/><circle cx="12" cy="10" r="2.5"/><path d="M12 7.5v5M9.5 10h5"/></svg></template>เพิ่มจุดบริการ</AppButton>
          <AppButton size="md" variant="outline" :disabled="projects.length === 0" @click="togglePolygonEdit"><template #prefix><svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M12 20h9"/><path d="M16.5 3.5a2.1 2.1 0 0 1 3 3L8 18l-4 1 1-4Z"/></svg></template>แก้ไขโซน</AppButton>
          <AppButton size="md" variant="danger-outline" :disabled="!activeProject && !overviewEditProject" @click="requestDeleteActiveProject"><template #prefix><svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M3 6h18M8 6V4h8v2M19 6l-1 14H6L5 6M10 11v5M14 11v5"/></svg></template>ลบโซน</AppButton>
        </template>

        <template v-else-if="isDrawingZone || isDrawingStall || drawingPoints.length > 0">
          <AppButton size="md" variant="danger-outline" @click="cancelDrawing">
            <template #prefix><svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="m6 6 12 12M18 6 6 18"/></svg></template>
            ยกเลิก
          </AppButton>
          <template v-if="isDrawingStall">
            <AppButton size="md" :variant="stallDrawShape === 'rectangle' ? 'primary' : 'outline'" @click="setStallDrawShape('rectangle')">
              <template #prefix><svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><rect x="3" y="5" width="18" height="14" rx="2"/></svg></template>
              ทรงสี่เหลี่ยม
            </AppButton>
            <AppButton size="md" :variant="stallDrawShape === 'polygon' ? 'primary' : 'outline'" @click="setStallDrawShape('polygon')">
              <template #prefix><svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><polygon points="9 2 15 2 19 6 19 12 15 16 9 16 5 12 5 6"/></svg></template>
              วาดเอง
            </AppButton>
            <AppButton
              v-if="stallDrawShape === 'rectangle'"
              size="md"
              :variant="stallRectAlign === 'zone' ? 'primary' : 'outline'"
              @click="toggleStallRectAlign"
            >
              <template #prefix><svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M4 20 20 4"/><rect x="7" y="7" width="10" height="10" rx="1.5" transform="rotate(-15 12 12)"/></svg></template>
              {{ stallRectAlign === 'zone' ? 'ขนานขอบโซน' : 'แนวตรง' }}
            </AppButton>
          </template>
          <AppButton size="md" variant="outline" :disabled="!drawingPoints.length && !rectStartCorner" @click="undoDrawingPoint">
            <template #prefix><svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="m15 18-6-6 6-6"/><path d="M9 12h10"/></svg></template>
            {{ isDrawingStall && stallDrawShape === 'rectangle' ? 'เริ่มกรอบใหม่' : 'ย้อนจุด' }}
          </AppButton>
          <AppButton size="md" variant="primary" :disabled="drawingPoints.length < 3" @click="finishDrawing">
            <template #prefix><svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.2"><path d="m5 12 4 4L19 6"/></svg></template>
            <span v-if="isDrawingStall && stallDrawShape === 'rectangle'">
              บันทึกแผงที่วาด{{ drawnStallDimensions ? ` (${drawnStallDimensions.widthM}×${drawnStallDimensions.depthM} ม.)` : '' }}
            </span>
            <span v-else>{{ isDrawingStall ? 'บันทึกแผงที่วาด' : 'บันทึกโซนที่วาด' }} ({{ drawingPoints.length }} จุด)</span>
          </AppButton>
        </template>

        <template v-else-if="isAddingFacility">
          <AppButton size="md" variant="danger-outline" @click="cancelDrawing">
            <template #prefix><svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="m6 6 12 12M18 6 6 18"/></svg></template>
            ยกเลิกการปักจุดบริการ
          </AppButton>
        </template>

        <template v-else-if="isDeletingStall">
          <AppButton size="md" variant="danger-outline" @click="cancelDrawing">
            <template #prefix><svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="m6 6 12 12M18 6 6 18"/></svg></template>
            ยกเลิกการลบแผง
          </AppButton>
        </template>

        <template v-else-if="isPolygonEditMode">
          <AppButton size="md" variant="danger-outline" @click="cancelPolygonEdit">
            <template #prefix><svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="m6 6 12 12M18 6 6 18"/></svg></template>
            ยกเลิก
          </AppButton>
          <template v-if="polygonEditTargetProject">
            <AppButton size="md" variant="outline" :disabled="!polygonEditHistory.length" @click="undoPolygonEdit">
              <template #prefix><svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="m15 18-6-6 6-6"/><path d="M9 12h10"/></svg></template>
              ย้อนกลับ
            </AppButton>
            <AppButton size="md" variant="danger-outline" @click="requestDeleteActiveProject">
              <template #prefix><svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M3 6h18M8 6V4h8v2M19 6l-1 14H6L5 6M10 11v5M14 11v5"/></svg></template>
              ลบโซน ({{ polygonEditTargetProject.name }})
            </AppButton>
            <AppButton size="md" variant="primary" @click="savePolygonEdit">
              <template #prefix><svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.2"><path d="m5 12 4 4L19 6"/></svg></template>
              บันทึกขอบเขตโซน ({{ polygonEditTargetProject.name }})
            </AppButton>
          </template>
        </template>

        <template v-else-if="isStallPolygonEditMode">
          <AppButton size="md" variant="danger-outline" @click="cancelStallPolygonEdit">
            <template #prefix><svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="m6 6 12 12M18 6 6 18"/></svg></template>
            ยกเลิก
          </AppButton>
          <AppButton size="md" variant="outline" :disabled="!stallEditHistory.length" @click="undoStallPolygonEdit">
            <template #prefix><svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="m15 18-6-6 6-6"/><path d="M9 12h10"/></svg></template>
            ย้อนกลับ
          </AppButton>
          <AppButton size="md" variant="danger-outline" :disabled="!editingStallCode" @click="deleteCurrentlyEditingStall">
            <template #prefix><svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M3 6h18M8 6V4h8v2M19 6l-1 14H6L5 6M10 11v5M14 11v5"/></svg></template>
            ลบแผง{{ editingStallCode ? ` (${editingStallCode})` : '' }}
          </AppButton>
          <AppButton size="md" variant="primary" @click="saveStallPolygonEdit">
            <template #prefix><svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.2"><path d="m5 12 4 4L19 6"/></svg></template>
            บันทึกขนาดแผง{{ editingStallCode ? ` (${editingStallCode})` : '' }}
          </AppButton>
        </template>
      </div>

      <div v-if="managementMode && isDeletingStall" class="polygon-edit-hint">
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><rect x="4" y="4" width="16" height="16" rx="2"/><path d="m9 9 6 6M15 9l-6 6"/></svg>
        <span>โหมดลบแผง: คลิกที่แผงเช่าบนแผนที่ที่ต้องการลบออกจากโซน</span>
      </div>

      <div v-if="managementMode && (isDrawingZone || isDrawingStall || drawingPoints.length > 0)" class="polygon-edit-hint">
        <svg v-if="isDrawingStall && stallDrawShape === 'rectangle'" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><rect x="3" y="5" width="18" height="14" rx="2"/></svg>
        <svg v-else viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><polygon points="9 2 15 2 19 6 19 12 15 16 9 16 5 12 5 6"/></svg>
        <template v-if="isDrawingStall && stallDrawShape === 'rectangle'">
          <span v-if="!rectStartCorner">โหมดทรงสี่เหลี่ยม ({{ stallRectAlign === 'zone' ? 'ขนานขอบโซน' : 'แนวตรง' }}): คลิกจุดที่ 1 บนแผนที่เพื่อกำหนดมุมเริ่มต้นของแผง</span>
          <span v-else-if="drawingPoints.length < 4">
            คลิกจุดที่ 2 เพื่อวางมุมตรงข้าม{{ liveRectPreview ? ` (ขนาดประมาณ ${liveRectPreview.widthM} × ${liveRectPreview.depthM} ม.)` : '' }}
          </span>
          <span v-else>
            ได้กรอบสี่เหลี่ยม{{ drawnStallDimensions ? ` ${drawnStallDimensions.widthM} × ${drawnStallDimensions.depthM} ม.` : '' }} แล้ว — กดปุ่ม “บันทึกแผงที่วาด” ด้านล่าง หรือคลิกบนแผนที่เพื่อวางกรอบใหม่
          </span>
        </template>
        <template v-else>
          <span v-if="drawingPoints.length < 3">คลิกบนแผนที่เพื่อวางจุดมุมอย่างน้อย 3 จุด (วางแล้ว {{ drawingPoints.length }}/3 จุด)</span>
          <span v-else>วางครบ {{ drawingPoints.length }} จุดแล้ว — กดปุ่ม “{{ isDrawingStall ? 'บันทึกแผงที่วาด' : 'บันทึกโซนที่วาด' }}” ด้านล่างเพื่อบันทึก</span>
        </template>
      </div>

      <div v-if="managementMode && isPolygonEditMode" class="polygon-edit-hint">
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M12 20h9"/><path d="M16.5 3.5a2.1 2.1 0 0 1 3 3L8 18l-4 1 1-4Z"/></svg>
        <span v-if="!polygonEditTargetProject">
          โหมดแก้ไขขอบเขตโซน: คลิกเลือกโซนหรือโครงการบนแผนที่ที่ต้องการปรับแก้ไขขอบเขต
        </span>
        <span v-else-if="!activeProject">
          กำลังแก้ไขขอบเขตโซน ({{ polygonEditTargetProject.name }}): ลากจุดมุมสีน้ำเงินเพื่อขยาย • คลิกหรือลากปุ่ม + บนเส้นขอบเพื่อเพิ่มมุมใหม่ (คลิกโซนอื่นเพื่อสลับโซน)
        </span>
        <span v-else>
          กำลังแก้ไขขอบเขตโซน ({{ polygonEditTargetProject.name }}): ลากจุดมุมสีน้ำเงินเพื่อปรับขอบเขต • คลิกหรือลากปุ่ม + บนเส้นขอบเพื่อเพิ่มมุมใหม่ (ดับเบิลคลิกจุดมุมเพื่อลบ)
        </span>
      </div>

      <div v-if="managementMode && isStallPolygonEditMode" class="polygon-edit-hint">
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><rect x="4" y="4" width="16" height="16" rx="2"/><circle cx="4" cy="4" r="2"/><circle cx="20" cy="4" r="2"/><circle cx="20" cy="20" r="2"/><circle cx="4" cy="20" r="2"/></svg>
        <span>โหมดปรับขนาดแผง ({{ editingStallCode || 'คลิกเลือกแผง' }}): ลากจุดมุมสีเขียวเพื่อปรับขนาด หรือคลิก/ลากปุ่ม + บนเส้นขอบเพื่อเพิ่มมุมใหม่</span>
      </div>

      <!-- Floating Map Top Controls: แผนที่ vs ดาวเทียม & ปุ่มปรับกึ่งกลาง (มุมซ้ายบน) -->
      <div class="map-floating-top-controls">
        <div
          v-if="marketFloors.length > 1 && !activeProject"
          class="market-floor-switcher"
          aria-label="เลือกระดับชั้นของตลาด"
        >
          <button
            type="button"
            class="market-floor-tab"
            :class="{ active: selectedMarketFloorId === 'all' }"
            title="แสดงโครงการจากทุกชั้น"
            @click="selectMarketFloor('all')"
          >
            ทั้งหมด
          </button>
          <button
            v-for="floor in marketFloors"
            :key="floor.id"
            type="button"
            class="market-floor-tab"
            :class="{ active: selectedMarketFloorId === floor.id }"
            :title="`${getFloorDisplayLabel(floor)} · ${floor.projectCount} โครงการ · ${floor.stallCount.toLocaleString()} แผง`"
            @click="selectMarketFloor(floor.id)"
          >
            {{ getFloorDisplayLabel(floor) }}
          </button>
        </div>

        <div
          v-else-if="activeProject && (activeProject.floors?.length || 0) > 1"
          class="market-floor-switcher active-project-floor-tabs"
          aria-label="เลือกระดับชั้นของโซน"
        >
          <button
            v-for="floor in activeProject.floors"
            :key="floor.id"
            type="button"
            class="market-floor-tab"
            :class="{ active: String(currentFloorId) === String(floor.id) }"
            :title="`${getFloorDisplayLabel(floor)} · ทั้งหมด ${getFloorLayout(activeProject, floor.id).total} แผง (ว่าง ${getFloorLayout(activeProject, floor.id).available} แผง)`"
            @click="switchFloor(floor.id)"
          >
            {{ getFloorDisplayLabel(floor) }} ({{ getFloorLayout(activeProject, floor.id).total }})
          </button>
        </div>

        <div v-if="!isOverviewPage" class="map-floating-style-switcher">
          <button
            type="button"
            class="map-mode-tab"
            :class="{ active: currentStyle === 'standard' }"
            @click="setMapStyle('standard')"
            title="มุมมองแผนที่ (Street Map)"
          >
            แผนที่
          </button>
          <button
            type="button"
            class="map-mode-tab"
            :class="{ active: currentStyle === 'satellite' }"
            @click="setMapStyle('satellite')"
            title="มุมมองภาพถ่ายดาวเทียม (Satellite Imagery)"
          >
            ดาวเทียม
          </button>
        </div>

        <button
          v-if="!isOverviewPage"
          type="button"
          class="map-floating-recenter-btn"
          @click="recenterMapView"
          title="ปรับมุมมองกลับสู่กึ่งกลาง"
        >
          <svg width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.2" stroke-linecap="round" stroke-linejoin="round">
            <circle cx="12" cy="10" r="10"></circle>
            <line x1="22" y1="12" x2="18" y2="12"></line>
            <line x1="6" y1="12" x2="2" y2="12"></line>
            <line x1="12" y1="6" x2="12" y2="2"></line>
            <line x1="12" y1="22" x2="12" y2="18"></line>
          </svg>
          <span>ปรับกึ่งกลาง</span>
        </button>
      </div>

      <!-- Floating Floor Preview Pill in Step 3 Preview Mode -->
      <transition name="fade">
        <div
          v-if="previewMode && activeProject && activeProjectFloors.length > 1"
          class="map-floating-floor-preview-pill"
        >
          <span class="floor-preview-label">ชั้น:</span>
          <div class="floor-preview-tabs">
            <button
              v-for="fl in activeProjectFloors"
              :key="fl.id"
              type="button"
              class="floor-preview-tab"
              :class="{ active: currentFloorId === fl.id }"
              @click="switchFloor(fl.id)"
            >
              {{ getFloorDisplayLabel(fl) }}
            </button>
          </div>
        </div>
      </transition>

      <!-- Floating Return to Booked Stalls Button (แสดงเมื่อคลิกดูโครงการอื่นใน Step 3) -->
      <!-- Return to booked project floating CTA (ทั้ง Step 1 และ Step 3) -->
      <transition name="fade">
        <button
          v-if="isNotViewingBookedProject && (selectedStallCodes.length > 0 || previewMode)"
          type="button"
          class="map-floating-return-btn"
          @click="focusBookedStalls"
          title="คลิกเพื่อกลับไปยังตำแหน่งแผงที่คุณเลือกทันที"
        >
          <span class="pulse-indicator-dot"></span>
          <svg width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round">
            <path d="M21 10c0 7-9 13-9 13s-9-6-9-13a9 9 0 0 1 18 0z"></path>
            <circle cx="12" cy="10" r="3"></circle>
          </svg>
          <span style="font-weight: 700;">กลับไปยังแผงที่คุณเลือก ({{ bookedProject?.name || 'โครงการที่เลือก' }})</span>
          <span class="floating-stall-codes-badge">{{ selectedStallCodes.join(', ') }}</span>
        </button>
      </transition>

      <!-- Overlay slot: ให้หน้าจองวาง drawer แผงที่เลือกทับบนแผนที่ -->
      <slot name="overlay" :activeProject="activeProject" />

      <!-- Multi-floor projects require an explicit floor choice before showing stalls. -->
      <transition name="fsc-fade">
        <div
          v-if="floorSelectorProject"
          class="floor-selector-modal-overlay"
          role="dialog"
          aria-modal="true"
          aria-labelledby="floor-selector-title"
          @click.self="closeFloorSelector"
        >
          <section class="floor-selector-panel">
            <AppModalHeader
              title="เลือกชั้นที่ต้องการดู"
              :subtitle="`${floorSelectorProject.name} · ${floorSelectorProject.floors?.length} ชั้น · โครงการที่ ${floorSelectorProject.sectionNumber}`"
              icon="grid"
              @close="closeFloorSelector"
            />

            <div class="fsc-floors-stack" aria-label="รายการชั้น">
              <button
                v-for="floor in floorSelectorProject.floors"
                :key="floor.id"
                type="button"
                class="fsc-floor-item"
                :class="{
                  'is-floor-full': getFloorLayout(floorSelectorProject, floor.id).available === 0 && getFloorLayout(floorSelectorProject, floor.id).reserved === 0,
                  'is-floor-queue': getFloorLayout(floorSelectorProject, floor.id).available === 0 && getFloorLayout(floorSelectorProject, floor.id).reserved > 0
                }"
                @click="selectFloorFromDialog(floor.id)"
              >
                <span class="fsc-pillar" :class="{ 'is-full': getFloorLayout(floorSelectorProject, floor.id).available === 0 && getFloorLayout(floorSelectorProject, floor.id).reserved === 0 }">
                  <span class="fsc-pillar-level">{{ getFloorDisplayLabel(floor) }}</span>
                  <span class="fsc-pillar-badge">ระดับ {{ floor.level }}</span>
                </span>
                <span class="fsc-floor-details">
                  <span class="fsc-floor-headline">
                    <span class="fsc-floor-title-text">{{ getFloorDisplayLabel(floor) }}</span>
                    <span v-if="getFloorLayout(floorSelectorProject, floor.id).available > 0" class="fsc-avail-chip">
                      <i class="fsc-chip-dot"></i>ว่าง {{ getFloorLayout(floorSelectorProject, floor.id).available }}
                    </span>
                    <span v-else-if="getFloorLayout(floorSelectorProject, floor.id).reserved > 0" class="fsc-avail-chip is-queue">
                      <i class="fsc-chip-dot"></i>คิวสำรอง {{ getFloorLayout(floorSelectorProject, floor.id).reserved }}
                    </span>
                    <span v-else class="fsc-avail-chip is-full">
                      <i class="fsc-chip-dot"></i>เต็มแล้ว
                    </span>
                  </span>
                  <span class="fsc-floor-meta-row">
                    <span class="fsc-meta-pill">ทั้งหมด {{ getFloorLayout(floorSelectorProject, floor.id).total }} แผง</span>
                    <span v-if="getFloorLayout(floorSelectorProject, floor.id).reserved > 0" class="fsc-meta-pill">ติดจอง (ยังไม่ทำสัญญา) {{ getFloorLayout(floorSelectorProject, floor.id).reserved }} แผง</span>
                    <span v-if="getFloorLayout(floorSelectorProject, floor.id).available === 0 && getFloorLayout(floorSelectorProject, floor.id).reserved === 0" class="fsc-meta-pill is-full-tag">มีผู้เช่าครบ 100%</span>
                  </span>
                </span>
                <span class="fsc-action-col">
                  <span v-if="getFloorLayout(floorSelectorProject, floor.id).available > 0" class="fsc-select-action-btn">
                    ดูผังแผง <span>→</span>
                  </span>
                  <span v-else-if="getFloorLayout(floorSelectorProject, floor.id).reserved > 0" class="fsc-select-action-btn is-queue-action">
                    จองคิวสำรอง <span>→</span>
                  </span>
                  <span v-else class="fsc-select-action-btn is-full-action">
                    ดูผัง (เต็ม) <span>→</span>
                  </span>
                </span>
              </button>
            </div>
          </section>
        </div>
      </transition>

      <!-- Drill-down Panel: แสดงเมื่อคลิกเข้าโครงการ (แผงย่อยวาดอยู่ในกรอบบนแผนที่) -->
      <transition name="smp-slide">
        <aside
          v-if="!previewMode && activeProject && activeLayout && !appliedSearchQuery.trim()"
          class="stall-mode-panel"
          :class="{
            'is-swapped-out': hideProjectPanel,
            'is-expanding-to-stall': expandProjectPanel,
            'is-returning-from-stall': returnProjectPanel,
            'is-exiting-project': isProjectPanelExiting,
            'is-sheet-expanded': isProjectSheetExpanded
          }"
          :style="projectSheetHeight !== null ? {
            height: `${projectSheetHeight}px !important`,
            maxHeight: `${projectSheetHeight}px !important`,
            transition: isSheetDragging ? 'none !important' : 'height 0.22s cubic-bezier(0.16, 1, 0.3, 1), max-height 0.22s cubic-bezier(0.16, 1, 0.3, 1) !important'
          } : undefined"
          :aria-hidden="hideProjectPanel"
          aria-label="รายละเอียดโครงการ"
        >
          <!-- ขีดสำหรับลากปรับขนาดหรือรูดปิด Bottom Sheet บนโหมด Responsive -->
          <div
            class="smp-drag-handle-bar"
            role="button"
            tabindex="0"
            aria-label="ปรับขนาดหรือรูดปิดแผงข้อมูลโครงการ"
            @touchstart="onSheetTouchStart"
            @touchmove="onSheetTouchMove"
            @touchend="onSheetTouchEnd"
          >
            <span class="smp-drag-handle-pill"></span>
          </div>

          <!-- เนื้อหาภายในเลื่อน Scroll ได้อิสระ โดยไม่ดันหรือทับขีดด้านบน -->
          <div class="smp-scroll-content">
            <div class="smp-header-row smp-desktop-only">
            <div class="smp-back-section">
              <button id="btn-back-overview" type="button" class="smp-back-btn" @click="handleBackOverviewClick">
                <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"><polyline points="15 18 9 12 15 6"></polyline></svg>
                <span>ดูผังรวมตลาด</span>
              </button>
            </div>
            <button
              type="button"
              class="smp-expand-toggle-btn"
              :title="isProjectSheetExpanded ? 'ย่อลง' : 'ขยายขนาดแผง'"
              aria-label="ย่อหรือขยายแผงโครงการ"
              @click="toggleProjectSheetExpand"
            >
              <svg v-if="!isProjectSheetExpanded" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.2" stroke-linecap="round" stroke-linejoin="round">
                <polyline points="18 15 12 9 6 15"></polyline>
              </svg>
              <svg v-else width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.2" stroke-linecap="round" stroke-linejoin="round">
                <polyline points="6 9 12 15 18 9"></polyline>
              </svg>
            </button>
          </div>

          <div class="smp-title-row">
            <span class="smp-dot" :style="{ backgroundColor: activeProject.color }"></span>
            <h3 class="smp-title">{{ activeProject.name }}</h3>
          </div>
          <div class="smp-category">{{ activeProject.category }}</div>

          <!-- Multi-floor Selector: เรียบง่าย กดสลับชั้นได้ทันทีในแถบด้านข้าง ไม่บังหน้าจอ -->
          <div v-if="activeProjectFloors.length > 1" class="smp-floor-section">
            <div class="smp-floor-header">
              <div class="smp-floor-title">
                <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.2" stroke-linecap="round" stroke-linejoin="round">
                  <path d="M2 20h20"></path>
                  <path d="M5 20V4h14v16"></path>
                  <path d="M9 8h2M9 12h2M9 16h2M13 8h2M13 12h2M13 16h2"></path>
                </svg>
                <span>ระดับชั้นอาคาร</span>
              </div>
              <span class="smp-floor-count">{{ activeProjectFloors.length }} ชั้น</span>
            </div>
            <div
              class="smp-floor-segmented-row"
              :class="{
                'has-many-floors': activeProjectFloors.length >= 3,
                'has-five-or-more-floors': activeProjectFloors.length >= 5
              }"
            >
              <button
                v-for="fl in activeProjectFloors"
                :key="fl.id"
                type="button"
                class="smp-floor-segmented-btn"
                :class="{
                  active: currentFloorId === fl.id,
                  'is-floor-full': getFloorAvailCount(fl.id) === 0 && getFloorReservedCount(fl.id) === 0,
                  'is-floor-queue': getFloorAvailCount(fl.id) === 0 && getFloorReservedCount(fl.id) > 0
                }"
                @click="switchFloor(fl.id)"
              >
                <span class="smp-seg-badge">{{ fl.level || fl.id }}</span>
                <div class="smp-seg-info">
                  <span class="smp-seg-name">{{ getFloorDisplayLabel(fl) }}</span>
                  <span v-if="getFloorAvailCount(fl.id) > 0" class="smp-seg-avail">ว่าง {{ getFloorAvailCount(fl.id) }}</span>
                  <span v-else-if="getFloorReservedCount(fl.id) > 0" class="smp-seg-avail is-queue">คิวสำรอง {{ getFloorReservedCount(fl.id) }}</span>
                  <span v-else class="smp-seg-avail is-full">เต็ม</span>
                </div>
              </button>
            </div>
          </div>

          <!-- เนื้อหาแผงโครงการ: หน้าข้อมูลตลาดบางมุมมอง (เช่น จดมิเตอร์) ส่งเนื้อหาของตัวเองผ่านสล็อตแทนสถิติแผงเดิม -->
          <slot v-if="$slots['project-body']" name="project-body" :project="activeProject" />
          <template v-else>
          <!-- Full Floor Alert Banner with Quick Switch to Available Floor -->
          <div v-if="activeProjectFloors.length > 1 && activeLayout?.available === 0 && activeLayout?.reserved === 0" class="smp-floor-full-banner">
            <div class="sffb-head">
              <svg width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.2" stroke-linecap="round" stroke-linejoin="round">
                <circle cx="12" cy="12" r="10"></circle>
                <line x1="12" y1="16" x2="12" y2="12"></line>
                <line x1="12" y1="8" x2="12.01" y2="8"></line>
              </svg>
              <span>{{ currentFloorObject?.name || 'ชั้นนี้' }} มีผู้เช่าครบแล้ว</span>
            </div>
            <p class="sffb-desc">ชั้นนี้มีผู้เช่าทำสัญญาครบทุกแผงแล้ว ท่านสามารถดูผังหรือกดสลับไปชั้นที่มีแผงว่างได้ทันที</p>
            <button
              v-if="suggestedAvailableFloor"
              type="button"
              class="sffb-switch-btn"
              @click="switchFloor(suggestedAvailableFloor.id)"
            >
              <span>สลับไป {{ getFloorDisplayLabel(suggestedAvailableFloor) }} (ว่าง {{ getFloorAvailCount(suggestedAvailableFloor.id) }} แผง)</span>
              <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"><polyline points="9 18 15 12 9 6"></polyline></svg>
            </button>
          </div>

          <!-- Queue / Reserved Floor Banner -->
          <div v-else-if="activeProjectFloors.length > 1 && activeLayout?.available === 0 && (activeLayout?.reserved || 0) > 0" class="smp-floor-queue-banner">
            <div class="sfqb-head">
              <svg width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.2" stroke-linecap="round" stroke-linejoin="round">
                <circle cx="12" cy="12" r="10"></circle>
                <polyline points="12 6 12 12 14 14"></polyline>
              </svg>
              <span>{{ currentFloorObject?.name || 'ชั้นนี้' }} ติดจอง (เปิดรับคิวสำรอง)</span>
            </div>
            <p class="sfqb-desc">มีแผงติดจองที่ยังไม่ทำสัญญาจำนวน {{ activeLayout?.reserved }} แผง ท่านสามารถคลิกเลือกแผงสีส้มเพื่อลงทะเบียนคิวสำรอง (Waiting List) ได้ทันที</p>
            <button
              v-if="suggestedAvailableFloor"
              type="button"
              class="sfqb-switch-btn"
              @click="switchFloor(suggestedAvailableFloor.id)"
            >
              <span>หรือเลือกดูแผงว่างที่ {{ getFloorDisplayLabel(suggestedAvailableFloor) }} (ว่าง {{ getFloorAvailCount(suggestedAvailableFloor.id) }} แผง)</span>
              <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"><polyline points="9 18 15 12 9 6"></polyline></svg>
            </button>
          </div>

          <!-- 4 สถิติแผงค้าหลัก (แผงทั้งหมด / ว่าง / ติดจอง / มีผู้เช่า): หน้าข้อมูลตลาดมีการ์ด KPI ด้านบนแสดงตัวเลขชุดเดียวกันอยู่แล้ว จึงซ่อนไม่ให้ซ้ำ -->
          <div v-if="!isOverviewPage" class="smp-stats" :class="{ 'is-loading': statusDataLoading }" :aria-busy="statusDataLoading">
            <button type="button" class="smp-stat is-total" :class="{ 'is-active': appliedAvailabilityFilter === 'all' }" :aria-pressed="appliedAvailabilityFilter === 'all'" :disabled="statusDataLoading" @click="setProjectStatusQuickFilter('all')">
              <span class="smp-stat-label">แผงทั้งหมด</span>
              <span v-if="statusDataLoading" class="smp-stat-skeleton" aria-hidden="true"></span>
              <strong v-else>{{ activeLayout.total.toLocaleString('th-TH') }}</strong>
            </button>
            <button type="button" class="smp-stat is-available" :class="{ 'is-active': appliedAvailabilityFilter === 'available' }" :aria-pressed="appliedAvailabilityFilter === 'available'" :disabled="statusDataLoading" @click="setProjectStatusQuickFilter('available')">
              <span class="smp-stat-label">ว่างพร้อมเช่า</span>
              <span v-if="statusDataLoading" class="smp-stat-skeleton" aria-hidden="true"></span>
              <strong v-else>{{ activeLayout.available.toLocaleString('th-TH') }}</strong>
            </button>
            <button type="button" class="smp-stat is-reserved" :class="{ 'is-active': appliedAvailabilityFilter === 'reserved' }" :aria-pressed="appliedAvailabilityFilter === 'reserved'" :disabled="statusDataLoading" @click="setProjectStatusQuickFilter('reserved')">
              <span class="smp-stat-label">ติดจอง</span>
              <span v-if="statusDataLoading" class="smp-stat-skeleton" aria-hidden="true"></span>
              <strong v-else>{{ activeLayout.reserved.toLocaleString('th-TH') }}</strong>
            </button>
            <button type="button" class="smp-stat is-occupied" :class="{ 'is-active': appliedAvailabilityFilter === 'occupied' }" :aria-pressed="appliedAvailabilityFilter === 'occupied'" :disabled="statusDataLoading" @click="setProjectStatusQuickFilter('occupied')">
              <span class="smp-stat-label">มีผู้เช่า</span>
              <span v-if="statusDataLoading" class="smp-stat-skeleton" aria-hidden="true"></span>
              <strong v-else>{{ activeLayout.occupied.toLocaleString('th-TH') }}</strong>
            </button>
          </div>

          <!-- อัตราการเช่าพื้นที่ (Occupancy Progress Bar) -->
          <div v-if="managementMode" class="smp-occ-card" :class="{ 'is-map-colors': isOverviewPage }">
            <div class="smp-occ-header">
              <div class="smp-occ-title-group">
                <span class="smp-occ-title">อัตราการเช่าพื้นที่</span>
                <small class="smp-occ-sub">สัดส่วนผู้เช่าและติดจอง</small>
              </div>
              <strong class="smp-occ-pct">
                {{ activeLayout.total > 0 ? (((activeLayout.occupied + activeLayout.reserved) / activeLayout.total) * 100).toFixed(1) : '0.0' }}%
              </strong>
            </div>
            <div class="smp-occ-track" role="progressbar" :aria-valuenow="activeLayout.total > 0 ? (((activeLayout.occupied + activeLayout.reserved) / activeLayout.total) * 100).toFixed(0) : 0" aria-valuemin="0" aria-valuemax="100">
              <span class="track-occupied" :style="{ width: `${activeLayout.total > 0 ? (activeLayout.occupied / activeLayout.total) * 100 : 0}%` }"></span>
              <span class="track-reserved" :style="{ width: `${activeLayout.total > 0 ? (activeLayout.reserved / activeLayout.total) * 100 : 0}%` }"></span>
            </div>
            <div class="smp-occ-legend">
              <span><i class="dot-occ"></i> {{ isOverviewPage ? 'มีผู้เช่า' : 'เช่าแล้ว' }} {{ activeLayout.total > 0 ? ((activeLayout.occupied / activeLayout.total) * 100).toFixed(0) : 0 }}%</span>
              <span><i class="dot-res"></i> ติดจอง {{ activeLayout.total > 0 ? ((activeLayout.reserved / activeLayout.total) * 100).toFixed(0) : 0 }}%</span>
              <span><i class="dot-avail"></i> ว่าง {{ activeLayout.total > 0 ? ((activeLayout.available / activeLayout.total) * 100).toFixed(0) : 0 }}%</span>
            </div>
          </div>

          <!-- ข้อมูลโครงการ: ขนาดพื้นที่ & ขนาดแผงมาตรฐาน -->
          <div class="smp-project-info-grid" :class="{ 'is-list': isOverviewPage || !managementMode }">
            <div class="smp-info-pill">
              <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                <rect x="3" y="3" width="18" height="18" rx="2"></rect><path d="M3 9h18M9 21V9"></path>
              </svg>
              <div class="smp-info-text">
                <small>พื้นที่ในโครงการ</small>
                <strong class="smp-area-measure"><span>{{ activeProject.realAreaSqm ? activeProject.realAreaSqm.toLocaleString('th-TH') : '-' }}</span><em>ตร.ม.</em></strong>
              </div>
            </div>
            <div class="smp-info-pill">
              <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                <rect x="2" y="7" width="20" height="14" rx="2"></rect><path d="M16 21V5a2 2 0 0 0-2-2h-4a2 2 0 0 0-2 2v16"></path>
              </svg>
              <div class="smp-info-text">
                <small>ขนาดแผงมาตรฐาน</small>
                <strong><span>2.5 × 2.0</span><em>ม. (5 ตร.ม.)</em></strong>
              </div>
            </div>
          </div>

          <!-- สัญลักษณ์สีสถานะแผงบนผัง -->
          <div v-if="managementMode && !isOverviewPage" class="smp-legend">
            <span class="smp-legend-title">สัญลักษณ์สีแผงค้า</span>
            <div class="smp-legend-items">
              <div><span class="smp-sw sw-available"></span>ว่างพร้อมเช่า</div>
              <div><span class="smp-sw sw-reserved"></span>ติดจอง / คิวสำรอง</div>
              <div><span class="smp-sw sw-occupied"></span>มีผู้เช่าแล้ว</div>
            </div>
          </div>

          <div class="smp-hint" :class="{ 'is-quiet': isOverviewPage }">
            <svg width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="10"></circle><line x1="12" y1="16" x2="12" y2="12"></line><line x1="12" y1="8" x2="12.01" y2="8"></line></svg>
            <span v-if="managementMode && inspectOnly">{{ isOverviewPage ? 'แตะแผงบนผังเพื่อดูรายละเอียดแผง ผู้เช่า และสัญญา' : 'แตะแผงใดก็ได้บนผังแผนที่เพื่อดูข้อมูลแผง ผู้เช่า และสัญญาเช่า' }}</span>
            <span v-else-if="managementMode">แตะแผงบนผังแผนที่เพื่อตั้งค่าข้อมูลแผง หรือใช้ปุ่มเครื่องมือด้านบนเพื่อวาดและปรับขนาด</span>
            <span v-else>แตะแผงสีเขียวหรือสีส้มบนผังแผนที่เพื่อเลือกแผงค้า</span>
          </div>
          </template>
        </div>
      </aside>
      </transition>

      <!-- First-visit guidance: disappears as soon as a project is selected. -->
      <transition name="fade">
        <div v-if="!previewMode && !managementMode && !activeProject" class="map-project-guidance" aria-live="polite">
          <span class="map-project-guidance-icon">
            <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.2" stroke-linecap="round" stroke-linejoin="round"><path d="M3 21h18"></path><path d="M5 21V5l7-3 7 3v16"></path><path d="M9 9h1M14 9h1M9 13h1M14 13h1"></path></svg>
          </span>
          <span>
            <strong>เลือกโครงการที่ต้องการจอง</strong>
            <small>คลิกพื้นที่กรอบสีบนแผนที่เพื่อดูแผงว่าง</small>
          </span>
          <svg class="map-project-guidance-arrow" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.2" stroke-linecap="round" stroke-linejoin="round"><path d="M5 12h14"></path><path d="m13 6 6 6-6 6"></path></svg>
        </div>
      </transition>

<!-- หน้าข้อมูลตลาด: แถบเครื่องมือ 4 ปุ่ม (แดชบอร์ด / ค้นหา / กรอง / ชั้นข้อมูล) + แผงสไลด์จากขวา -->
      <template v-if="isOverviewPage && !previewMode">
        <nav
          class="ov-rail"
          :class="{
            'is-shifted': !!ovPanel && ovPanel !== 'dashboard',
            'is-shifted-dashboard': ovPanel === 'dashboard',
            'is-hidden': hideRail,
            'is-project-active': !!activeProject
          }"
          aria-label="เครื่องมือแผนที่"
        >
          <button
            type="button"
            class="ov-rail-btn"
            :class="{ active: ovPanel === 'dashboard' }"
            title="แดชบอร์ดสรุปภาพรวม"
            aria-label="แดชบอร์ดสรุปภาพรวม"
            :aria-expanded="ovPanel === 'dashboard'"
            aria-controls="ov-slide-panel"
            @click="toggleOvPanel('dashboard')"
          >
            <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true">
              <rect x="3" y="3" width="7" height="7" rx="1"></rect>
              <rect x="14" y="3" width="7" height="7" rx="1"></rect>
              <rect x="14" y="14" width="7" height="7" rx="1"></rect>
              <rect x="3" y="14" width="7" height="7" rx="1"></rect>
            </svg>
          </button>
          <button type="button" class="ov-rail-btn" :class="{ active: ovPanel === 'search' }" title="ค้นหา" aria-label="ค้นหา" :aria-expanded="ovPanel === 'search'" aria-controls="ov-slide-panel" @click="toggleOvPanel('search')">
            <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true"><circle cx="11" cy="11" r="8"></circle><line x1="21" y1="21" x2="16.65" y2="16.65"></line></svg>
          </button>
          <button type="button" class="ov-rail-btn" :class="{ active: ovPanel === 'filter' }" title="ตัวกรอง" aria-label="ตัวกรอง" :aria-expanded="ovPanel === 'filter'" aria-controls="ov-slide-panel" @click="toggleOvPanel('filter')">
            <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true"><path d="M4 6h16M7 12h10m-7 6h4"></path></svg>
            <span v-if="ovFilterBadge" class="ov-rail-badge">{{ ovFilterBadge }}</span>
          </button>
          <button type="button" class="ov-rail-btn" :class="{ active: ovPanel === 'layers' }" title="ชั้นข้อมูล" aria-label="ชั้นข้อมูล" :aria-expanded="ovPanel === 'layers'" aria-controls="ov-slide-panel" @click="toggleOvPanel('layers')">
            <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true"><polygon points="12 2 2 7 12 12 22 7 12 2"></polygon><polyline points="2 17 12 22 22 17"></polyline><polyline points="2 12 12 17 22 12"></polyline></svg>
            <span v-if="mapLayerFilterCount" class="ov-rail-badge">{{ mapLayerFilterCount }}</span>
          </button>
        </nav>

        <transition name="ov-slide">
          <aside
            v-if="ovPanel"
            id="ov-slide-panel"
            class="ov-slide"
            :class="{ 'is-dashboard-panel': ovPanel === 'dashboard' }"
            role="dialog"
            aria-modal="false"
            :aria-label="ovTitle"
          >
            <header class="ov-slide-head">
              <span class="ov-slide-icon" aria-hidden="true">
                <svg v-if="ovPanel === 'dashboard'" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                  <rect x="3" y="3" width="7" height="7" rx="1"></rect>
                  <rect x="14" y="3" width="7" height="7" rx="1"></rect>
                  <rect x="14" y="14" width="7" height="7" rx="1"></rect>
                  <rect x="3" y="14" width="7" height="7" rx="1"></rect>
                </svg>
                <svg v-else-if="ovPanel === 'search'" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="11" cy="11" r="8"></circle><line x1="21" y1="21" x2="16.65" y2="16.65"></line></svg>
                <svg v-else-if="ovPanel === 'filter'" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M4 6h16M7 12h10m-7 6h4"></path></svg>
                <svg v-else width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><polygon points="12 2 2 7 12 12 22 7 12 2"></polygon><polyline points="2 17 12 22 22 17"></polyline><polyline points="2 12 12 17 22 12"></polyline></svg>
              </span>
              <div class="ov-slide-title">
                <h3>{{ ovTitle }}</h3>
                <p>{{ ovSubtitle }}</p>
              </div>
              <button type="button" class="ov-slide-close" aria-label="ปิดแผง" @click="closeOvPanel">
                <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="m6 6 12 12M18 6 6 18"></path></svg>
              </button>
            </header>

            <div class="ov-slide-body">
              <!-- แดชบอร์ดสรุปภาพรวม -->
              <template v-if="ovPanel === 'dashboard'">
                <slot name="dashboard" />
              </template>

              <!-- ค้นหา -->
              <template v-else-if="ovPanel === 'search'">
                <label class="ov-search">
                  <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true"><circle cx="11" cy="11" r="8"></circle><line x1="21" y1="21" x2="16.65" y2="16.65"></line></svg>
                  <input
                    ref="ovSearchInput"
                    v-model="searchQuery"
                    type="search"
                    :placeholder="managementMode ? 'เลขที่การจอง คำร้อง ชื่อโครงการ หรือรหัสแผง' : 'ชื่อโครงการหรือประเภทสินค้า'"
                    @input="handleSearchInput"
                    @keydown.down.prevent="moveSearchSuggestion(1)"
                    @keydown.up.prevent="moveSearchSuggestion(-1)"
                    @keydown.enter.prevent="selectActiveSearchSuggestion"
                  />
                  <button v-if="searchQuery" type="button" class="ov-search-clear" aria-label="ล้างคำค้นหา" @click="clearSearch">
                    <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="m6 6 12 12M18 6 6 18"></path></svg>
                  </button>
                </label>
                <div v-if="appliedSearchQuery.trim()" class="ov-applied-search">
                  <span>กำลังแสดงผล: <strong>{{ appliedSearchQuery }}</strong></span>
                  <button type="button" @click="clearAppliedSearch">ล้างผลบนแผนที่</button>
                </div>
                <AppEmptyState
                  v-if="!suggestionQuery.trim() || !searchResults.length"
                  variant="map"
                  icon="search"
                  :title="suggestionQuery.trim() ? 'ไม่พบผลลัพธ์' : 'พิมพ์คำค้นหา'"
                  :description="suggestionQuery.trim() ? `ไม่พบ “${suggestionQuery}” ลองใช้คำอื่น` : 'ใช้เลขที่การจอง เลขที่คำร้อง ชื่อโครงการ หรือรหัสแผง'"
                />
                <template v-else>
                  <p class="ov-search-summary">
                    พบ <strong>{{ searchResults.length.toLocaleString('th-TH') }}</strong> โครงการ
                    <template v-if="searchResultStallCount"> · <strong>{{ searchResultStallCount.toLocaleString('th-TH') }}</strong> แผง</template>
                  </p>
                  <ul class="ov-result-list" role="listbox" aria-label="ผลการค้นหา">
                  <li v-for="(result, index) in searchResults" :key="result.project.id">
                    <button
                      type="button"
                      class="ov-result"
                      :class="{
                        focused: index === focusedSuggestionIndex,
                        active: String(result.project.id) === selectedSuggestionId
                      }"
                      :data-suggestion-index="index"
                      role="option"
                      :aria-selected="String(result.project.id) === selectedSuggestionId"
                      @click="selectSearchResult(result)"
                    >
                      <span class="ov-swatch" :style="{ backgroundColor: projectColor(result.project) }"></span>
                      <span class="ov-result-text">
                        <strong>
                          <template v-for="(part, partIndex) in getHighlightedSearchParts(result.project.name)" :key="partIndex">
                            <mark v-if="part.match">{{ part.text }}</mark><template v-else>{{ part.text }}</template>
                          </template>
                        </strong>
                        <span class="ov-result-badges">
                          <span v-for="type in result.matchTypes" :key="type" class="ov-result-badge" :class="`is-${type}`">{{ searchMatchTypeLabel(type) }}</span>
                        </span>
                        <small v-if="result.referenceText" class="ov-result-reference">
                          <template v-for="(part, partIndex) in getHighlightedSearchParts(result.referenceText)" :key="partIndex">
                            <mark v-if="part.match">{{ part.text }}</mark><template v-else>{{ part.text }}</template>
                          </template>
                        </small>
                        <small>{{ result.project.category || 'ไม่ระบุหมวด' }} · {{ (result.project.stalls?.length || result.project.totalStalls || 0).toLocaleString('th-TH') }} แผง</small>
                      </span>
                      <svg class="ov-chevron" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true"><polyline points="9 18 15 12 9 6"></polyline></svg>
                    </button>
                  </li>
                  </ul>
                </template>
              </template>

              <!-- ตัวกรอง: ชุดตัวกรองต่างกันตามโหมด -->
              <template v-else-if="ovPanel === 'filter'">
                <slot v-if="filterContext === 'costs'" name="filter-costs" />
                <template v-else>
                  <div v-if="activeKpiFilterLabel" class="ov-kpi-filter-chip">
                    <span>สถานะจากการ์ด: <strong>{{ activeKpiFilterLabel }}</strong></span>
                    <button type="button" :aria-label="`ล้างสถานะ ${activeKpiFilterLabel}`" @click="emit('clear-status-filter')">
                      <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true"><path d="m6 6 12 12M18 6 6 18"></path></svg>
                    </button>
                  </div>

                  <div v-if="contractStatusFilter !== 'all'" class="ov-kpi-filter-chip">
                    <span>สถานะสัญญา: <strong>{{ contractStatusLabel }}</strong></span>
                    <button type="button" :aria-label="`ล้างสถานะสัญญา ${contractStatusLabel}`" @click="applyDashboardFilter({ contractStatus: 'all' })">
                      <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true"><path d="m6 6 12 12M18 6 6 18"></path></svg>
                    </button>
                  </div>

                  <section class="ov-filter-section">
                    <h4>พื้นที่</h4>
                    <div class="ov-filter-grid">
                      <div class="ov-field">
                        <span class="ov-label">โครงการ / โซน</span>
                        <AppSelect v-model="zoneFilter" :options="zoneOptions" placeholder="ทั้งหมด" />
                      </div>
                      <div class="ov-field">
                        <span class="ov-label">ชั้น</span>
                        <AppSelect v-model="floorFilter" :options="floorFilterOptions" placeholder="ทุกชั้น" />
                      </div>
                    </div>
                  </section>

                  <section class="ov-filter-section">
                    <h4>คุณสมบัติแผง</h4>
                    <div class="ov-filter-grid">
                      <div class="ov-field">
                        <span class="ov-label">ค่าเช่ารายเดือน</span>
                        <AppSelect v-model="rentFilter" :options="rentOptions" placeholder="ทุกช่วงราคา" />
                      </div>
                      <div class="ov-field">
                        <span class="ov-label">ขนาดแผง</span>
                        <AppSelect v-model="sizeFilter" :options="sizeOptions" placeholder="ทุกขนาด" />
                      </div>
                      <div class="ov-field">
                        <span class="ov-label">ความครบถ้วนของข้อมูล</span>
                        <AppSelect v-model="readinessFilter" :options="readinessOptions" placeholder="ทั้งหมด" />
                      </div>
                    </div>
                  </section>

                  <section class="ov-filter-section">
                    <h4>ช่วงเวลา</h4>
                    <div class="ov-filter-grid">
                      <div class="ov-field">
                        <span class="ov-label">สัญญาจะหมดอายุ</span>
                        <AppSelect v-model="contractExpiryFilter" :options="contractExpiryOptions" placeholder="ไม่จำกัด" />
                        <div v-if="contractExpiryFilter === 'custom'" class="ov-date-range">
                          <label><span>ตั้งแต่</span><AppDatePicker v-model="contractExpiryFrom" placeholder="เลือกวันที่" aria-label="สัญญาหมดอายุตั้งแต่วันที่" /></label>
                          <label><span>ถึง</span><AppDatePicker v-model="contractExpiryTo" placeholder="เลือกวันที่" align="right" aria-label="สัญญาหมดอายุถึงวันที่" /></label>
                        </div>
                      </div>
                      <div class="ov-field">
                        <span class="ov-label">มีการจองในช่วง</span>
                        <AppSelect v-model="bookingRangeFilter" :options="bookingRangeOptions" placeholder="ไม่จำกัด" />
                        <div v-if="bookingRangeFilter === 'custom'" class="ov-date-range">
                          <label><span>ตั้งแต่</span><AppDatePicker v-model="bookingFrom" placeholder="เลือกวันที่" aria-label="มีการจองตั้งแต่วันที่" /></label>
                          <label><span>ถึง</span><AppDatePicker v-model="bookingTo" placeholder="เลือกวันที่" align="right" aria-label="มีการจองถึงวันที่" /></label>
                        </div>
                      </div>
                    </div>
                  </section>

                  <section class="ov-filter-section">
                    <h4>สัญญาและการส่งคืน</h4>
                    <div class="ov-filter-grid">
                      <div class="ov-field">
                        <span class="ov-label">การต่ออายุสัญญา</span>
                        <AppSelect v-model="renewalStatusFilter" :options="renewalStatusOptions" placeholder="ทุกสถานะ" />
                      </div>
                      <div class="ov-field">
                        <span class="ov-label">การส่งคืนแผง</span>
                        <AppSelect v-model="returnStatusFilter" :options="returnStatusOptions" placeholder="ทุกสถานะ" />
                      </div>
                    </div>
                  </section>

                  <section v-if="statusFilterContext === 'legal_case'" class="ov-filter-section">
                    <h4>รายละเอียดคดี</h4>
                    <div class="ov-filter-grid">
                      <div class="ov-field">
                        <span class="ov-label">ขั้นตอนของคดี</span>
                        <AppSelect v-model="legalStatusFilter" :options="legalStatusOptions" placeholder="ทุกขั้นตอน" />
                      </div>
                    </div>
                  </section>

                  <section v-if="statusFilterContext === 'action_required'" class="ov-filter-section">
                    <h4>รายละเอียดที่ต้องดำเนินการ</h4>
                    <div class="ov-filter-grid">
                      <div class="ov-field">
                        <span class="ov-label">สถานะงานซ่อม</span>
                        <AppSelect v-model="maintenanceStatusFilter" :options="maintenanceStatusFilterOptions" placeholder="ทุกสถานะ" />
                      </div>
                      <div class="ov-field">
                        <span class="ov-label">ประเภทคำร้อง</span>
                        <AppSelect v-model="requestTypeFilter" :options="requestTypeOptions" placeholder="ทุกประเภท" />
                      </div>
                    </div>
                  </section>
                  <p class="ov-result-summary">ผลลัพธ์ <strong>{{ previewFilteredProjectCount }}</strong> โครงการ · <strong>{{ previewFilteredStallCount.toLocaleString('th-TH') }}</strong> แผง</p>
                </template>
              </template>

              <!-- ชั้นข้อมูล -->
              <template v-else>
                <section class="ov-section">
                  <h4>แผนที่ฐาน</h4>
                  <div class="ov-basemaps">
                    <button v-for="b in basemapOptions" :key="b.value" type="button" class="ov-basemap" :class="{ active: currentStyle === b.value }" :aria-pressed="currentStyle === b.value" @click="setMapStyle(b.value)">
                      <span class="ov-basemap-thumb" :class="`is-${b.value}`" :data-keep-colors="b.value === 'satellite' ? '' : null" aria-hidden="true">
                        <svg v-if="currentStyle === b.value" class="ov-basemap-check" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3" stroke-linecap="round" stroke-linejoin="round"><polyline points="20 6 9 17 4 12"></polyline></svg>
                      </span>
                      <span>{{ b.label }}</span>
                    </button>
                  </div>
                </section>
                <section class="ov-section">
                  <div class="ov-section-head">
                    <h4>สิ่งอำนวยความสะดวก</h4>
                  </div>
                  <div class="ov-layer-toggle-list" role="group" aria-label="สิ่งอำนวยความสะดวก">
                    <button
                      type="button"
                      class="ov-layer-toggle"
                      :class="{ active: showFacilities }"
                      :aria-pressed="showFacilities"
                      @click="toggleFacilitiesLayer"
                    >
                      <span class="ov-layer-icon" aria-hidden="true">
                        <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                          <path d="M3 20h18"></path>
                          <path d="M5 20v-2a7 7 0 0 1 14 0v2"></path>
                          <path d="M12 11V4"></path>
                          <path d="M10 4h4"></path>
                        </svg>
                      </span>
                      <span class="ov-layer-text">
                        <strong>สิ่งอำนวยความสะดวก</strong>
                        <small>จุดอำนวยความสะดวก {{ managedFacilities.length.toLocaleString('th-TH') }} รายการ</small>
                      </span>
                      <span class="ov-layer-switch" aria-hidden="true">
                        <span></span>
                      </span>
                    </button>
                  </div>
                </section>
                <section class="ov-section">
                  <div class="ov-section-head">
                    <h4>ป้ายกำกับ</h4>
                  </div>
                  <div class="ov-label-layer-list" role="group" aria-label="ป้ายกำกับบนแผนที่">
                    <button type="button" class="ov-label-layer" :class="{ active: showProjectLabels }" :aria-pressed="showProjectLabels" @click="toggleProjectLabels">
                      <span class="ov-label-layer-icon" aria-hidden="true">
                        <svg width="17" height="17" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M4 5h16v14H4z"></path><path d="M8 9h8M8 13h5"></path></svg>
                      </span>
                      <span>ชื่อโครงการ</span>
                      <span class="ov-mini-switch" aria-hidden="true"><span></span></span>
                    </button>
                    <button type="button" class="ov-label-layer" :class="{ active: showStallNumberLabels }" :aria-pressed="showStallNumberLabels" @click="toggleStallLabels">
                      <span class="ov-label-layer-icon" aria-hidden="true">
                        <svg width="17" height="17" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M4 4h6v6H4zM14 4h6v6h-6zM4 14h6v6H4zM14 14h6v6h-6z"></path></svg>
                      </span>
                      <span>เลขแผง</span>
                      <span class="ov-mini-switch" aria-hidden="true"><span></span></span>
                    </button>
                  </div>
                </section>
                <section class="ov-section">
                  <div class="ov-section-head">
                    <h4>สีโครงการ</h4>
                    <div class="ov-legend-actions">
                      <span class="ov-selected-count">{{ legendCategories.length ? `${legendCategories.length} ที่เลือก` : 'แสดงทั้งหมด' }}</span>
                      <button v-if="legendCategories.length" type="button" class="ov-show-all" @click="clearLegendCategories">แสดงทั้งหมด</button>
                    </div>
                  </div>
                  <ul class="ov-legend" role="group" aria-label="สีโครงการ">
                    <li v-for="cat in categoryStats" :key="cat.name">
                      <button type="button" class="ov-legend-item" :class="{ active: legendCategories.includes(cat.name), dimmed: legendCategories.length > 0 && !legendCategories.includes(cat.name) }" :style="legendCategoryStyle(cat.color)" :aria-pressed="legendCategories.includes(cat.name)" @click="toggleLegendCategory(cat.name)">
                        <span class="ov-tile" :style="{ backgroundColor: cat.color }" :title="`${cat.projectCount} โครงการ`">{{ cat.projectCount }}</span>
                        <span class="ov-legend-text">
                          <strong>{{ cat.name }}</strong>
                          <small>โครงการ {{ cat.projectNumbers.join(', ') }}</small>
                        </span>
                        <span class="ov-check" aria-hidden="true">
                          <svg v-if="legendCategories.includes(cat.name)" width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3.2" stroke-linecap="round" stroke-linejoin="round"><polyline points="20 6 9 17 4 12"></polyline></svg>
                        </span>
                      </button>
                    </li>
                  </ul>
                </section>
              </template>
            </div>

            <footer v-if="ovPanel === 'filter'" class="ov-slide-foot">
              <template v-if="filterContext === 'costs'">
                <AppButton variant="secondary" size="md" full-width @click="resetCostFiltersAndOverview">ล้างค่า</AppButton>
              </template>
              <template v-else>
                <AppButton variant="secondary" size="md" @click="clearAllFilters">ล้างค่า</AppButton>
                <AppButton variant="primary" size="md" @click="applyFiltersFromPanel">ดูผลลัพธ์ ({{ previewFilteredProjectCount }})</AppButton>
              </template>
            </footer>
          </aside>
        </transition>
      </template>

    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch, nextTick, onMounted, onBeforeUnmount, onActivated, onDeactivated } from 'vue'
import AppButton from '~/components/AppButton.vue'
import { addStandardMapControls } from '~/utils/mapControls'
import AppDatePicker from '~/components/AppDatePicker.vue'
import AppSelect, { type SelectOption } from '~/components/AppSelect.vue'
import type { ProjectZone, ProjectFloor, MarketFacility } from '../../data/market_gis_types'
import { getProjectStallLayout, getAllFloorLayoutsForProject, isOverviewBoundaryProject, clearLayoutCache, type MapStall, type ProjectStallLayout } from '../../data/market_stall_layout'
import { getMarketGisConfig, type MarketGisConfig } from '../../data/markets_registry'
import { loadMarketGisJson, mergeMarketProjects, cloneProjects, getLoadedMarketFacilities, sharedRequest } from '../../data/market_loader'

export type { ProjectZone, ProjectFloor }

export interface StallOperationalFilterMeta {
  contractStatus?: 'active' | 'expiring' | 'expired'
  contractExpiryDays?: number
  lastBookingAt?: number
  renewalStatus?: 'pending' | 'reviewing' | 'approved' | 'declined'
  returnStatus?: 'notified' | 'pending' | 'inspection' | 'returned' | 'overdue'
  legalStatus?: 'legal_review' | 'proceeding' | 'enforcement' | 'closed'
  maintenanceStatus?: 'pending' | 'repairing' | 'completed' | 'overdue'
  requestTypes?: Array<'contract_name_change' | 'transfer' | 'contract_cancel' | 'sublease'>
  searchTerms?: string[]
  searchReferences?: Array<{ type: 'booking' | 'application'; number: string }>
}

const props = withDefaults(defineProps<{
  containerId?: string
  marketId?: string | number
  marketFull?: boolean
  selectedProjectId?: string
  selectedStallCodes?: string[]
  selectedStallKeys?: string[]
  myBookedStallCodes?: string[]
  stallStatusOverrides?: Record<string, 'available' | 'reserved' | 'occupied'>
  stallStatusFilter?: 'all' | 'available' | 'reserved' | 'occupied' | 'action_required'
  statusFilterContext?: 'all' | 'available' | 'reserved' | 'occupied' | 'action_required' | 'legal_case'
  actionRequiredStallCodes?: string[]
  onlyAvailable?: boolean
  categoryFilter?: string | null
  mapStyle?: 'standard' | 'dark' | 'satellite'
  previewMode?: boolean
  managementMode?: boolean
  inspectOnly?: boolean
  isOverviewPage?: boolean
  filterContext?: 'data' | 'costs'
  filterBadge?: number
  hideRail?: boolean
  hideProjectPanel?: boolean
  expandProjectPanel?: boolean
  returnProjectPanel?: boolean
  managedProjects?: ProjectZone[]
  managedFacilities?: MarketFacility[]
  managedDataVersion?: number
  isDashboardOpen?: boolean
  dashboardTitle?: string
  dashboardSubtitle?: string
  activeFilterCodes?: string[] | null
  operationalFilterMeta?: Record<string, StallOperationalFilterMeta>
  statusDataLoading?: boolean
}>(), {
  containerId: 'step1-projects-maplibre',
  marketId: '1',
  marketFull: false,
  selectedProjectId: undefined,
  selectedStallCodes: () => [],
  selectedStallKeys: () => [],
  myBookedStallCodes: () => [],
  stallStatusOverrides: () => ({}),
  stallStatusFilter: 'all',
  statusFilterContext: 'all',
  actionRequiredStallCodes: () => [],
  onlyAvailable: false,
  categoryFilter: null,
  mapStyle: 'standard',
  previewMode: false,
  managementMode: false,
  inspectOnly: false,
  isOverviewPage: false,
  filterContext: 'data',
  filterBadge: 0,
  hideRail: false,
  hideProjectPanel: false,
  expandProjectPanel: false,
  returnProjectPanel: false,
  managedProjects: () => [],
  managedFacilities: () => [],
  managedDataVersion: 0,
  isDashboardOpen: false,
  dashboardTitle: '',
  dashboardSubtitle: '',
  activeFilterCodes: null,
  operationalFilterMeta: () => ({}),
  statusDataLoading: false
})

const emit = defineEmits<{
  (e: 'selectProject', project: ProjectZone): void
  (e: 'selectStall', stall: MapStall): void
  (e: 'updateProjectPolygon', project: ProjectZone): void
  (e: 'updateStallPolygon', project: ProjectZone, stall: any): void
  (e: 'requestAddProject', coordinates: [number, number][]): void
  (e: 'requestAddStallLayout', project: ProjectZone, coordinates: [number, number][], floorId: string, dimensions?: { widthM: number; depthM: number }): void
  (e: 'requestAddFacility', coordinates: [number, number], projectId?: string): void
  (e: 'requestDeleteFacility', facility: MarketFacility): void
  (e: 'requestDeleteStall', project: ProjectZone, stall: any): void
  (e: 'requestDeleteProject', project: ProjectZone, stats?: { total: number; available: number; reserved: number; occupied: number }): void
  (e: 'requestExitProject'): void
  (e: 'exitProject'): void
  (e: 'update:categoryFilter', cat: string | null): void
  (e: 'update:mapStyle', style: 'standard' | 'satellite' | 'dark'): void
  (e: 'reset-cost-filters'): void
  (e: 'clearStallSelection'): void
  (e: 'update:onlyAvailable', val: boolean): void
  (e: 'update:isDashboardOpen', val: boolean): void
  (e: 'toggleDashboard'): void
  (e: 'update:activePanel', val: string | null): void
  (e: 'clear-status-filter'): void
  (e: 'operational-filter-change', filters: { zone: string; search: string; contractStatus: string; renewalStatus: string; returnStatus: string; legalStatus: string; maintenanceStatus: string }): void
  (e: 'market-data-loaded', payload: { projects: ProjectZone[]; facilities: MarketFacility[] }): void
}>()

// Active Map Style Mode ('standard' = Street, 'satellite' = Imagery)
const currentStyle = ref<'standard' | 'satellite' | 'dark'>(props.mapStyle === 'satellite' ? 'satellite' : props.mapStyle === 'dark' ? 'dark' : 'standard')

// เมื่อระบบอยู่โหมดมืด แผนที่ฐานแบบ "แผนที่ปกติ" จะใช้สไตล์มืดโดยอัตโนมัติ (ดาวเทียมคงเดิม)
const { isDark: isSystemDark } = useTheme()
const { showWarning: showMapWarning } = useAppAlert()
const resolveStyle = (mode: 'standard' | 'satellite' | 'dark') => (mode === 'standard' && isSystemDark.value ? 'dark' : mode)
const setMapStyle = (mode: 'standard' | 'satellite' | 'dark') => {
  if (currentStyle.value === mode) return
  currentStyle.value = mode
  emit('update:mapStyle', mode)
  if (!map) return

  map.setStyle(getStyleConfig(resolveStyle(mode)))

  map.once('styledata', () => {
    setupProjectLayers()
    refreshSources()
  })
}

const searchQuery = ref('')
const suggestionQuery = ref('')
const focusedSuggestionIndex = ref(-1)
const selectedSuggestionId = ref<string | null>(null)
let searchDebounceTimer: ReturnType<typeof setTimeout> | null = null
const rentFilter = ref<'all' | 'under-2000' | '2000-3000' | 'over-3000'>('all')
type StallSizeFilter = 'all' | 'up-to-5' | '5-to-10' | 'over-10'
const sizeFilter = ref<StallSizeFilter>('all')
const availabilityFilter = ref<'available' | 'reserved' | 'occupied' | 'action_required' | 'all'>('all')
const zoneFilter = ref('all')
const readinessFilter = ref<'all' | 'ready' | 'missing-stalls' | 'incomplete-stalls'>('all')
type ContractStatusFilter = 'all' | 'active' | 'expiring' | 'expired'
type ExpiryFilter = 'all' | '30' | '60' | '90' | 'custom'
type BookingRangeFilter = 'all' | '7' | '30' | '90' | 'custom'
type RenewalFilter = 'all' | 'awaiting' | 'pending' | 'reviewing' | 'approved' | 'declined'
type ReturnFilter = 'all' | 'notified' | 'pending' | 'inspection' | 'returned' | 'overdue'
type LegalFilter = 'all' | 'legal_review' | 'proceeding' | 'enforcement' | 'closed'
type MaintenanceFilter = 'all' | 'pending' | 'repairing' | 'completed' | 'overdue'
type RequestTypeFilter = 'all' | 'contract_name_change' | 'transfer' | 'contract_cancel' | 'sublease'
const floorFilter = ref('all')
const contractStatusFilter = ref<ContractStatusFilter>('all')
const contractExpiryFilter = ref<ExpiryFilter>('all')
const contractExpiryFrom = ref('')
const contractExpiryTo = ref('')
const bookingRangeFilter = ref<BookingRangeFilter>('all')
const bookingFrom = ref('')
const bookingTo = ref('')
const renewalStatusFilter = ref<RenewalFilter>('all')
const returnStatusFilter = ref<ReturnFilter>('all')
const legalStatusFilter = ref<LegalFilter>('all')
const maintenanceStatusFilter = ref<MaintenanceFilter>('all')
const requestTypeFilter = ref<RequestTypeFilter>('all')
// Inputs keep draft values. The map only reads the applied values after the
// user presses “ค้นหา”, so changing a dropdown never shifts the map mid-task.
const draftCategoryFilter = ref(props.categoryFilter || '')
const appliedCategoryFilter = ref(props.categoryFilter || '')
// กรองหลายหมวดจากสีโครงการ (ชั้นข้อมูล) แยกจากตัวกรองหมวดเดี่ยวในหน้าตัวกรอง
const legendCategories = ref<string[]>([])
const showFacilities = ref(false)
const showStallNumberLabels = ref(true)
const showProjectLabels = ref(true)
const appliedSearchQuery = ref('')
const appliedRentFilter = ref<'all' | 'under-2000' | '2000-3000' | 'over-3000'>('all')
const appliedSizeFilter = ref<StallSizeFilter>('all')
const appliedAvailabilityFilter = ref<'available' | 'reserved' | 'occupied' | 'action_required' | 'all'>('all')
const appliedZoneFilter = ref('all')
const appliedReadinessFilter = ref<'all' | 'ready' | 'missing-stalls' | 'incomplete-stalls'>('all')
const appliedFloorFilter = ref('all')
const appliedContractStatusFilter = ref<ContractStatusFilter>('all')
const appliedContractExpiryFilter = ref<ExpiryFilter>('all')
const appliedContractExpiryFrom = ref('')
const appliedContractExpiryTo = ref('')
const appliedBookingRangeFilter = ref<BookingRangeFilter>('all')
const appliedBookingFrom = ref('')
const appliedBookingTo = ref('')
const appliedRenewalStatusFilter = ref<RenewalFilter>('all')
const appliedReturnStatusFilter = ref<ReturnFilter>('all')
const appliedLegalStatusFilter = ref<LegalFilter>('all')
const appliedMaintenanceStatusFilter = ref<MaintenanceFilter>('all')
const appliedRequestTypeFilter = ref<RequestTypeFilter>('all')

const setProjectStatusQuickFilter = (status: 'available' | 'reserved' | 'occupied' | 'all') => {
  const nextStatus = status !== 'all' && appliedAvailabilityFilter.value === status ? 'all' : status
  availabilityFilter.value = nextStatus
  appliedAvailabilityFilter.value = nextStatus
  emit('update:onlyAvailable', nextStatus === 'available')
  refreshSources()
}
const showAdvancedFilters = ref(false)
const showColorLegend = ref(false)
const activeKpiFilterLabel = computed(() => {
  const labels: Record<string, string> = {
    available: 'ว่าง',
    reserved: 'ติดจอง',
    occupied: 'มีผู้เช่า',
    action_required: 'ต้องดำเนินการ',
    legal_case: 'มีคดี'
  }
  return labels[props.statusFilterContext] || ''
})

const currentMarket = ref<MarketGisConfig>(getMarketGisConfig(props.marketId))
const managedFacilities = ref<MarketFacility[]>([])
const { get: apiGet } = useApi()
const currentMarketCenter = computed<[number, number]>(() => currentMarket.value.center)
const currentMarketZoom = computed<number>(() => currentMarket.value.overviewZoom)
const currentMarketBearing = computed<number>(() => currentMarket.value.bearing ?? DEFAULT_BEARING)
const currentMarketPitch = computed<number>(() => currentMarket.value.pitch ?? DEFAULT_PITCH)
const selectedMarketFloorId = ref('all')

type MarketFloorOption = ProjectFloor & {
  projectCount: number
  stallCount: number
}

// Prefer floor tags attached to the surveyed stalls. Project metadata can be
// broader than the actual GIS content (Min Buri originally declared both
// floors on every project even when its stalls existed on only one floor).
const getProjectActualFloorIds = (project: ProjectZone): string[] => {
  const declaredFloorIds = (project.floors || []).map(floor => String(floor.id))
  const stallFloorIds = Array.from(new Set(
    (project.stalls || [])
      .map(stall => stall.floor ?? stall.floorLevel)
      .filter(floor => floor !== undefined && floor !== null && String(floor).trim())
      .map(floor => String(floor))
  ))
  if (project.layoutMode === 'manual' || String(project.id).startsWith('zone-') || props.managementMode) {
    return Array.from(new Set([...declaredFloorIds, ...stallFloorIds]))
  }
  if (stallFloorIds.length) return stallFloorIds
  return declaredFloorIds
}

const marketFloors = computed<MarketFloorOption[]>(() => {
  const floors = new Map<string, MarketFloorOption>()
  currentMarket.value.projects.forEach(project => {
    getProjectActualFloorIds(project).forEach(floorId => {
      const metadata = project.floors?.find(floor => String(floor.id) === floorId)
      const existing = floors.get(floorId)
      const stallCount = (project.stalls || []).filter(stall => (
        String(stall.floor ?? stall.floorLevel ?? floorId) === floorId
      )).length
      if (existing) {
        existing.projectCount += 1
        existing.stallCount += stallCount
        return
      }
      const level = metadata?.level ?? (Number(floorId) || floors.size + 1)
      const thaiLevelLabel = `ชั้น ${level}`
      floors.set(floorId, {
        id: floorId,
        name: metadata?.name || thaiLevelLabel,
        shortLabel: metadata?.name || (metadata?.shortLabel && !/^\d+F$/i.test(metadata.shortLabel) ? metadata.shortLabel : thaiLevelLabel),
        level,
        projectCount: 1,
        stallCount
      })
    })
  })
  return Array.from(floors.values()).sort((a, b) => a.level - b.level)
})

const getFloorDisplayLabel = (floor?: MarketFloorOption | ProjectFloor | { name?: string; shortLabel?: string; level?: number; id?: string | number } | null | string): string => {
  if (!floor) return ''
  if (typeof floor === 'string') {
    const trimmed = floor.trim()
    if (trimmed === 'ชั้นล่าง') return 'ชั้น 1'
    if (trimmed === 'ชั้นบน') return 'ชั้น 2'
    if (/^(\d+)F$/i.test(trimmed)) return `ชั้น ${trimmed.replace(/F/i, '')}`
    if (/^B(\d+)$/i.test(trimmed)) return `ชั้นใต้ดิน B${trimmed.replace(/B/i, '')}`
    if (/^G$/i.test(trimmed) || /^GF$/i.test(trimmed)) return 'ชั้น G'
    if (trimmed.startsWith('ชั้น')) return trimmed
    if (/^\d+$/.test(trimmed)) return `ชั้น ${trimmed}`
    return trimmed
  }

  // 1. If name is available, clean and use it
  if (floor.name) {
    const rawName = floor.name.trim()
    if (rawName === 'ชั้นล่าง') return 'ชั้น 1'
    if (rawName === 'ชั้นบน') return 'ชั้น 2'
    if (rawName === 'ชั้นลอย') return rawName
    const matchFloor = rawName.match(/^(ชั้น\s*\d+|ชั้นล่าง|ชั้นบน|ชั้นลอย|ชั้น\s*[A-Za-z0-9]+)/)
    if (matchFloor) {
      return matchFloor[1]
    }
    if (rawName.startsWith('ชั้น')) {
      const firstPart = rawName.split(/[\(\/]/)[0]?.trim()
      if (firstPart) return firstPart
    }
    if (/^(\d+)F$/i.test(rawName)) return `ชั้น ${rawName.replace(/F/i, '')}`
    if (/^B(\d+)$/i.test(rawName)) return `ชั้นใต้ดิน B${rawName.replace(/B/i, '')}`
    if (/^G$/i.test(rawName) || /^GF$/i.test(rawName)) return 'ชั้น G'
    if (rawName) return rawName
  }

  // 2. Use level
  if (typeof floor.level === 'number' && !isNaN(floor.level)) {
    return `ชั้น ${floor.level}`
  }

  // 3. Fallback to shortLabel translated to Thai
  if (floor.shortLabel) {
    const sl = floor.shortLabel.trim()
    if (sl === 'ชั้นล่าง') return 'ชั้น 1'
    if (sl === 'ชั้นบน') return 'ชั้น 2'
    if (/^(\d+)F$/i.test(sl)) return `ชั้น ${sl.replace(/F/i, '')}`
    if (/^B(\d+)$/i.test(sl)) return `ชั้นใต้ดิน B${sl.replace(/B/i, '')}`
    if (/^G$/i.test(sl) || /^GF$/i.test(sl)) return 'ชั้น G'
    if (sl.startsWith('ชั้น')) return sl
    if (/^\d+$/.test(sl)) return `ชั้น ${sl}`
  }

  // 4. Fallback to id
  if (floor.id) {
    const idStr = String(floor.id).trim()
    if (/^\d+$/.test(idStr)) return `ชั้น ${idStr}`
    if (idStr.startsWith('ชั้น')) return idStr
    return `ชั้น ${idStr}`
  }

  return 'ชั้น 1'
}

type LngLat = [number, number]

// A few imported surveys also contain a coarse boundary around several real
// projects. Keep it in the data, but do not expose it as a selectable project.
const pointInProject = ([lng, lat]: LngLat, project: ProjectZone): boolean => {
  const ring = project.coordinates
  let inside = false
  for (let i = 0, j = ring.length - 1; i < ring.length; j = i++) {
    const [xi, yi] = ring[i]!
    const [xj, yj] = ring[j]!
    if ((yi > lat) !== (yj > lat) && lng < ((xj - xi) * (lat - yi)) / (yj - yi) + xi) inside = !inside
  }
  return inside
}

// This is used only to compare polygons from the same market, so a local
// planar area is both sufficient and more robust than their survey metadata.
const projectFootprintArea = (project: ProjectZone): number => {
  const ring = project.coordinates
  let twiceArea = 0
  for (let i = 0; i < ring.length - 1; i++) {
    const [x1, y1] = ring[i]!
    const [x2, y2] = ring[i + 1]!
    twiceArea += x1 * y2 - x2 * y1
  }
  return Math.abs(twiceArea) / 2
}

const baseMarketProjects = ref<ProjectZone[]>([])
let marketSyncRun = 0

const syncMarketJson = async (idOrSlug?: string | number, signal?: AbortSignal) => {
  const run = ++marketSyncRun
  const data = await loadMarketGisJson(idOrSlug, false, signal)
  if (data && run === marketSyncRun && !isUnmounted) {
    clearLayoutCache()
    clearProjectLayoutCache()
    // This data is used as a read-only map source. Keep the cached nested geometry by
    // reference and only replace the outer containers, otherwise a 10k-stall market
    // blocks the main thread while structuredClone duplicates the complete payload.
    // Read-only overview/booking maps can safely share the cached nested geometry.
    // Only the actual GIS editor needs an isolated mutable copy.
    const loadedProjects = props.managementMode && !props.isOverviewPage
      ? cloneProjects(data.projects || [])
      : (data.projects || [])
    baseMarketProjects.value = loadedProjects.filter(p => !String(p.id).startsWith('zone-'))
    currentMarket.value = { ...data, projects: loadedProjects }
    try {
      const marketKey = String(idOrSlug || data.id)
      const loadedFacilities = getLoadedMarketFacilities(marketKey)
      if (props.managementMode) {
        // The management page owns the draft request and pushes a versioned result
        // through managedProjects/managedFacilities. Do not request the same large
        // draft a second time from the child map.
        managedFacilities.value = props.managedFacilities.length
          ? props.managedFacilities
          : (loadedFacilities || [])
      } else if (loadedFacilities) {
        // gis-config ของตลาดนี้ถูกโหลดแล้วโดย loadMarketGisJson ไม่ต้องโหลดก้อนเดิมซ้ำ
        managedFacilities.value = loadedFacilities
      } else {
        const endpoint = `/markets/${encodeURIComponent(marketKey)}/gis-config`
        const response: any = await sharedRequest(endpoint, () => apiGet(endpoint))
        const managed = response?.data || response
        if (Array.isArray(managed?.projects) && managed.projects.length) {
          currentMarket.value = { ...data, projects: mergeMarketProjects(baseMarketProjects.value, managed.projects, false) }
        }
        managedFacilities.value = Array.isArray(managed?.facilities) ? managed.facilities : []
      }
    } catch {
      managedFacilities.value = []
    }
    if (props.managementMode && props.managedProjects.length) {
      currentMarket.value = { ...currentMarket.value, projects: mergeMarketProjects(baseMarketProjects.value, props.managedProjects, false) }
    }
    if (props.managementMode && props.managedFacilities.length) managedFacilities.value = props.managedFacilities
    emit('market-data-loaded', {
      projects: currentMarket.value.projects || [],
      facilities: managedFacilities.value
    })
    if (map) {
      refreshSources()
      renderZoneLabels()
    }
  }
}

const applyManagedProjects = (nextProjects: ProjectZone[]) => {
  if (!props.managementMode || !nextProjects.length) return
  clearLayoutCache()
  clearProjectLayoutCache()
  const merged = mergeMarketProjects(baseMarketProjects.value, nextProjects, false)
  currentMarket.value = { ...currentMarket.value, projects: merged }
  if (activeProjectId.value && !merged.some(project => project.id === activeProjectId.value)) {
    isPolygonEditMode.value = false
    clearPolygonVertexHandles()
    exitProject()
    return
  }
  if (props.selectedProjectId && props.selectedProjectId !== activeProjectId.value && !isPolygonEditMode.value) {
    const selectedTarget = projects.value.find(project => project.id === props.selectedProjectId)
    if (selectedTarget) {
      enterProject(selectedTarget)
      return
    }
  }
  refreshSources()
}

const applyManagedFacilities = (nextFacilities: MarketFacility[]) => {
  if (!props.managementMode) return
  managedFacilities.value = nextFacilities
  map?.getSource('market-facilities')?.setData(getFacilitiesGeoJson())
  renderFacilityMarkers()
}

const getFacilitiesGeoJson = () => ({
  type: 'FeatureCollection' as const,
  features: (showFacilities.value ? managedFacilities.value : [])
    .filter(facility => {
      const floorId = activeProjectId.value ? currentFloorId.value : selectedMarketFloorId.value
      if (floorId && floorId !== 'all' && facility.floorId && String(facility.floorId) !== String(floorId)) return false
      if (!activeProjectId.value) return true
      return findHostZoneForFacility(facility)?.id === activeProjectId.value
    })
    .map(facility => ({
      type: 'Feature' as const,
      geometry: facility.geometry || { type: 'Point' as const, coordinates: facility.coordinates },
      properties: {
        id: facility.id,
        name: facility.name,
        type: facility.type,
        floorId: facility.floorId || '',
        projectId: facility.projectId || ''
      }
    }))
})

// Extract categories dynamically from the current market's projects
const categories = computed(() => {
  const uniqueCats = new Map<string, string>()
  currentMarket.value.projects.forEach(p => {
    if (p.category && !uniqueCats.has(p.category)) {
      uniqueCats.set(p.category, p.color)
    }
  })
  return Array.from(uniqueCats.entries()).map(([name, color]) => ({ name, color }))
})

// Build normalized indexes once. The previous fallback walked every override for
// every stall, which became O(stalls × contracts) and froze large markets.
const stallStatusOverrideIndex = computed(() => {
  const direct = new Map<string, 'available' | 'reserved' | 'occupied'>()
  const digits = new Map<string, 'available' | 'reserved' | 'occupied'>()
  Object.entries(props.stallStatusOverrides || {}).forEach(([rawCode, status]) => {
    const cleanCode = rawCode.toUpperCase().trim()
    if (!cleanCode) return
    direct.set(cleanCode, status)
    const rawDigits = cleanCode.replace(/[^0-9]/g, '')
    if (!rawDigits) return
    if (!digits.has(rawDigits)) digits.set(rawDigits, status)
    const normalizedDigits = rawDigits.replace(/^0+/, '')
    if (normalizedDigits && !digits.has(normalizedDigits)) digits.set(normalizedDigits, status)
  })
  return { direct, digits }
})

// Helper to match stall status override with robust normalization
const findStallStatusOverride = (stallCode: string): 'available' | 'reserved' | 'occupied' | undefined => {
  const cleanCode = stallCode.toUpperCase().trim()
  const direct = stallStatusOverrideIndex.value.direct.get(cleanCode)
  if (direct) return direct
  // Fallback: match by number/identifier if project/zone is aligned (e.g. PRA-10-004 vs 10-A04 or PRA-10-4)
  const codeDigits = cleanCode.replace(/[^0-9]/g, '')
  if (!codeDigits) return undefined
  return stallStatusOverrideIndex.value.digits.get(codeDigits)
    || stallStatusOverrideIndex.value.digits.get(codeDigits.replace(/^0+/, ''))
}

// Apply the authoritative market-level state before any counter, color,
// filter, tooltip or click handler consumes a stall layout.
const applyMarketStatus = (layout: ProjectStallLayout): ProjectStallLayout => {
  if (!props.marketFull && Object.keys(props.stallStatusOverrides).length === 0) {
    return layout
  }
  const stalls = layout.stalls.map(stall => {
    const override = findStallStatusOverride(stall.code)
    return {
      ...stall,
      status: props.marketFull
        ? ('occupied' as const)
        : (override || stall.status || 'available')
    }
  })
  const available = stalls.filter(stall => stall.status === 'available').length
  const reserved = stalls.filter(stall => stall.status === 'reserved').length
  return {
    ...layout,
    stalls,
    available,
    reserved,
    occupied: stalls.length - available - reserved
  }
}

const getEffectiveFloorLayoutsForProject = (project: ProjectZone): Map<string, ProjectStallLayout> => {
  const effectiveLayouts = new Map<string, ProjectStallLayout>()
  getAllFloorLayoutsForProject(project).forEach((layout, floorId) => {
    effectiveLayouts.set(floorId, applyMarketStatus(layout))
  })
  return effectiveLayouts
}

// สร้าง layout เมื่อมีผู้ใช้เปิด/กรองโครงการนั้นเท่านั้น แทนการประมวลผล
// polygon ทุกแผงของทุกโครงการทุกครั้งที่ reactive state เปลี่ยน
const projectLayoutCache = new Map<string, ProjectStallLayout>()
const clearProjectLayoutCache = () => projectLayoutCache.clear()
const getProjectLayout = (project: ProjectZone): ProjectStallLayout => {
  const cached = projectLayoutCache.get(project.id)
  if (cached) return cached
  let result: ProjectStallLayout
  if (project.floors && project.floors.length > 1) {
    const allStalls: MapStall[] = []
    let total = 0
    let available = 0
    let reserved = 0
    let occupied = 0
    getAllFloorLayoutsForProject(project).forEach((floorLayout) => {
      const effective = applyMarketStatus(floorLayout)
      allStalls.push(...effective.stalls)
      total += effective.total
      available += effective.available
      reserved += effective.reserved
      occupied += effective.occupied
    })
    result = {
      stalls: allStalls,
      labelPoint: getProjectStallLayout(project).labelPoint,
      total,
      available,
      reserved,
      occupied,
      floorId: 'all',
      floorName: `รวมทุกชั้น (${project.floors.length} ชั้น)`
    }
  } else {
    result = applyMarketStatus(getProjectStallLayout(project))
  }
  projectLayoutCache.set(project.id, result)
  return result
}

const projects = computed<ProjectZone[]>(() => {
  const sourceProjects = currentMarket.value.projects
  return sourceProjects.filter(p => !isOverviewBoundaryProject(p, sourceProjects)).map(p => {
    const l = projectLayoutCache.get(p.id)
    return { ...p, totalStalls: l ? l.total : p.totalStalls, stallsAvailable: l ? l.available : p.stallsAvailable }
  })
})

const overviewProjects = computed<ProjectZone[]>(() => {
  if (marketFloors.value.length <= 1 || !selectedMarketFloorId.value || selectedMarketFloorId.value === 'all') return projects.value
  return projects.value.filter(project => getProjectActualFloorIds(project).includes(selectedMarketFloorId.value))
})

const selectMarketFloor = (floorId: string) => {
  if (selectedMarketFloorId.value === floorId) return
  selectedMarketFloorId.value = floorId
  hoveredProject.value = null
  hoverPopup?.remove()
  refreshSources()
  map?.flyTo({
    center: currentMarketCenter.value,
    zoom: currentMarketZoom.value,
    pitch: currentMarketPitch.value,
    bearing: currentMarketBearing.value,
    duration: 600
  })
}

const actionRequiredCodeSet = computed(() => new Set(
  (props.actionRequiredStallCodes || []).map(code => String(code || '').trim().toUpperCase()).filter(Boolean)
))

const activeFilterCodeSet = computed(() => {
  if (!props.activeFilterCodes) return null
  return new Set(props.activeFilterCodes.map(code => String(code || '').trim().toUpperCase()).filter(Boolean))
})

interface OperationalFilterSelection {
  floor: string
  contractStatus: ContractStatusFilter
  contractExpiry: ExpiryFilter
  contractExpiryFrom: string
  contractExpiryTo: string
  bookingRange: BookingRangeFilter
  bookingFrom: string
  bookingTo: string
  renewalStatus: RenewalFilter
  returnStatus: ReturnFilter
  legalStatus: LegalFilter
  maintenanceStatus: MaintenanceFilter
  requestType: RequestTypeFilter
}

const draftOperationalFilters = (): OperationalFilterSelection => ({
  floor: floorFilter.value,
  contractStatus: contractStatusFilter.value,
  contractExpiry: contractExpiryFilter.value,
  contractExpiryFrom: contractExpiryFrom.value,
  contractExpiryTo: contractExpiryTo.value,
  bookingRange: bookingRangeFilter.value,
  bookingFrom: bookingFrom.value,
  bookingTo: bookingTo.value,
  renewalStatus: renewalStatusFilter.value,
  returnStatus: returnStatusFilter.value,
  legalStatus: legalStatusFilter.value,
  maintenanceStatus: maintenanceStatusFilter.value,
  requestType: requestTypeFilter.value
})

const appliedOperationalFilters = (): OperationalFilterSelection => ({
  floor: appliedFloorFilter.value,
  contractStatus: appliedContractStatusFilter.value,
  contractExpiry: appliedContractExpiryFilter.value,
  contractExpiryFrom: appliedContractExpiryFrom.value,
  contractExpiryTo: appliedContractExpiryTo.value,
  bookingRange: appliedBookingRangeFilter.value,
  bookingFrom: appliedBookingFrom.value,
  bookingTo: appliedBookingTo.value,
  renewalStatus: appliedRenewalStatusFilter.value,
  returnStatus: appliedReturnStatusFilter.value,
  legalStatus: appliedLegalStatusFilter.value,
  maintenanceStatus: appliedMaintenanceStatusFilter.value,
  requestType: appliedRequestTypeFilter.value
})

// 'custom' ที่ยังไม่ได้เลือกวันที่ ถือว่ายังไม่กรอง (ช่องวันที่ว่าง = '')
const hasOperationalFilter = (filters: OperationalFilterSelection) => {
  const { contractExpiryFrom: expiryFrom, contractExpiryTo: expiryTo, bookingFrom: from, bookingTo: to, ...rest } = filters
  if (Object.values(rest).some(value => value !== 'all' && !(value === 'custom' && !(
    (rest.contractExpiry === 'custom' && (expiryFrom || expiryTo)) || (rest.bookingRange === 'custom' && (from || to))
  )))) return true
  return false
}

const dayStart = (dateText: string) => new Date(`${dateText}T00:00:00`).getTime()
const todayStart = () => { const d = new Date(); d.setHours(0, 0, 0, 0); return d.getTime() }

const matchesStallFiltersWith = (
  stall: MapStall,
  rent: typeof rentFilter.value,
  availability: typeof availabilityFilter.value,
  size: StallSizeFilter,
  operational: OperationalFilterSelection = appliedOperationalFilters()
) => {
  if (activeFilterCodeSet.value !== null) {
    const cleanCode = String(stall.code || '').trim().toUpperCase()
    if (!activeFilterCodeSet.value.has(cleanCode)) return false
  }
  if (availability === 'action_required') {
    const cleanCode = String(stall.code || '').trim().toUpperCase()
    if (!actionRequiredCodeSet.value.has(cleanCode)) return false
  } else if (availability !== 'all' && stall.status !== availability) {
    return false
  }
  if (rent === 'under-2000' && stall.monthlyRent >= 2000) return false
  if (rent === '2000-3000' && (stall.monthlyRent < 2000 || stall.monthlyRent > 3000)) return false
  if (rent === 'over-3000' && stall.monthlyRent <= 3000) return false
  if (size !== 'all') {
    const areaSqM = Number(stall.widthM) * Number(stall.depthM)
    if (!Number.isFinite(areaSqM) || areaSqM <= 0) return false
    if (size === 'up-to-5' && areaSqM > 5) return false
    if (size === '5-to-10' && (areaSqM <= 5 || areaSqM > 10)) return false
    if (size === 'over-10' && areaSqM <= 10) return false
  }
  const floorId = String(stall.floor ?? stall.floorLevel ?? '')
  if (operational.floor !== 'all' && floorId !== operational.floor) return false

  if (hasOperationalFilter({ ...operational, floor: 'all' })) {
    const meta = props.operationalFilterMeta[String(stall.code || '').trim().toUpperCase()]
    if (!meta) return false
    if (operational.contractStatus !== 'all' && meta.contractStatus !== operational.contractStatus) return false
    if (operational.contractExpiry === 'custom') {
      // ช่วงวันที่กำหนดเอง: เทียบวันหมดอายุกับช่วง (วันที่ผ่านมาแล้วก็เลือกได้ เช่น สัญญาที่หมดไปเมื่อเดือนก่อน)
      if (operational.contractExpiryFrom || operational.contractExpiryTo) {
        const expiryDays = Number(meta.contractExpiryDays)
        if (!Number.isFinite(expiryDays)) return false
        const base = todayStart()
        const fromDays = operational.contractExpiryFrom ? Math.round((dayStart(operational.contractExpiryFrom) - base) / 86_400_000) : -Infinity
        const toDays = operational.contractExpiryTo ? Math.round((dayStart(operational.contractExpiryTo) - base) / 86_400_000) : Infinity
        // เลือกวันเริ่มหลังวันสิ้นสุด: สลับให้เป็นช่วงที่ถูกต้อง
        const low = Math.min(fromDays, toDays)
        const high = Math.max(fromDays, toDays)
        if (expiryDays < low || expiryDays > high) return false
      }
    } else if (operational.contractExpiry !== 'all') {
      const expiryDays = Number(meta.contractExpiryDays)
      if (!Number.isFinite(expiryDays) || expiryDays < 0 || expiryDays > Number(operational.contractExpiry)) return false
    }
    if (operational.bookingRange !== 'all') {
      const bookedAt = Number(meta.lastBookingAt)
      if (operational.bookingRange === 'custom') {
        if (operational.bookingFrom || operational.bookingTo) {
          if (!Number.isFinite(bookedAt)) return false
          const startOf = operational.bookingFrom ? dayStart(operational.bookingFrom) : -Infinity
          const endOf = operational.bookingTo ? dayStart(operational.bookingTo) + 86_400_000 - 1 : Infinity
          // เลือกวันเริ่มหลังวันสิ้นสุด: สลับให้เป็นช่วงที่ถูกต้อง
          const from = startOf <= endOf ? startOf : endOf - 86_400_000 + 1
          const to = startOf <= endOf ? endOf : startOf + 86_400_000 - 1
          if (bookedAt < from || bookedAt > to) return false
        }
      } else if (!Number.isFinite(bookedAt) || bookedAt < Date.now() - Number(operational.bookingRange) * 86_400_000) return false
    }
    if (operational.renewalStatus === 'awaiting') {
      if (meta.renewalStatus !== 'pending' && meta.renewalStatus !== 'reviewing') return false
    } else if (operational.renewalStatus !== 'all' && meta.renewalStatus !== operational.renewalStatus) return false
    if (operational.returnStatus !== 'all' && meta.returnStatus !== operational.returnStatus) return false
    if (operational.legalStatus !== 'all' && meta.legalStatus !== operational.legalStatus) return false
    if (operational.maintenanceStatus !== 'all' && meta.maintenanceStatus !== operational.maintenanceStatus) return false
    if (operational.requestType !== 'all' && !meta.requestTypes?.includes(operational.requestType)) return false
  }
  return true
}

const matchesStallFilters = (stall: MapStall) => matchesStallFiltersWith(
  stall,
  appliedRentFilter.value,
  appliedAvailabilityFilter.value,
  appliedSizeFilter.value,
  appliedOperationalFilters()
)

const stallMatchesSearchWith = (stall: MapStall, value: string): boolean => {
  const query = value.trim().toLowerCase()
  if (!query) return true
  const code = String(stall.code || '').trim().toUpperCase()
  const mapNumber = getStallMapLabel(stall).toLowerCase()
  if (mapNumber.includes(query)) return true
  if (!/^\d+$/.test(query) && code.toLowerCase().includes(query)) return true
  return (props.operationalFilterMeta[code]?.searchTerms || [])
    .some(term => String(term).toLowerCase().includes(query))
}

const projectHasMatchingStallWith = (
  project: ProjectZone,
  rent: typeof rentFilter.value,
  availability: typeof availabilityFilter.value,
  size: StallSizeFilter,
  operational: OperationalFilterSelection
) => {
  // With no stall-level filter every project already matches. Avoid walking
  // and cloning thousands of stall polygons on every project hover.
  if (rent === 'all' && availability === 'all' && size === 'all' && !hasOperationalFilter(operational) && activeFilterCodeSet.value === null) return true
  const projectLayout = getProjectLayout(project)
  return Boolean(projectLayout?.stalls.some(stall => matchesStallFiltersWith(stall, rent, availability, size, operational)))
}

const getProjectReadiness = (project: ProjectZone): 'ready' | 'missing-stalls' | 'incomplete-stalls' => {
  const stalls = getProjectLayout(project).stalls
  if (!stalls.length) return 'missing-stalls'
  return stalls.some(stall => !stall.code || !Number(stall.widthM) || !Number(stall.depthM) || !Number(stall.monthlyRent))
    ? 'incomplete-stalls'
    : 'ready'
}

const projectMatchesAllFiltersWith = (
  project: ProjectZone,
  category: string,
  query: string,
  rent: typeof rentFilter.value,
  availability: typeof availabilityFilter.value,
  size: StallSizeFilter,
  zone = 'all',
  readiness: typeof readinessFilter.value = 'all',
  operational: OperationalFilterSelection = appliedOperationalFilters()
) => {
  const matchesCategory = (!category || project.category === category) &&
    (legendCategories.value.length === 0 || legendCategories.value.includes(project.category))
  const matchesZone = zone === 'all' || project.id === zone
  const matchesReadiness = readiness === 'all' || getProjectReadiness(project) === readiness
  return matchesCategory && matchesZone && matchesReadiness && matchesSearchWith(project, query) && projectHasMatchingStallWith(project, rent, availability, size, operational)
}

const projectMatchesAllFilters = (project: ProjectZone) => projectMatchesAllFiltersWith(
  project,
  appliedCategoryFilter.value,
  appliedSearchQuery.value,
  appliedRentFilter.value,
  appliedAvailabilityFilter.value,
  appliedSizeFilter.value,
  appliedZoneFilter.value,
  appliedReadinessFilter.value,
  appliedOperationalFilters()
)

const filteredProjectCount = computed(() => overviewProjects.value.filter(projectMatchesAllFilters).length)
const filteredStallCount = computed(() => overviewProjects.value.reduce((count, project) => {
  if (!projectMatchesAllFilters(project)) return count
  if (appliedRentFilter.value === 'all' && appliedAvailabilityFilter.value === 'all' && appliedSizeFilter.value === 'all' && !hasOperationalFilter(appliedOperationalFilters())) {
    return count + (project.stalls?.length || project.totalStalls || 0)
  }
  return count + getProjectLayout(project).stalls.filter(matchesStallFilters).length
}, 0))

const previewFilteredProjectCount = computed(() => overviewProjects.value.filter(project => projectMatchesAllFiltersWith(
  project,
  draftCategoryFilter.value,
  searchQuery.value,
  rentFilter.value,
  availabilityFilter.value,
  sizeFilter.value,
  zoneFilter.value,
  readinessFilter.value,
  draftOperationalFilters()
)).length)
const previewFilteredStallCount = computed(() => overviewProjects.value.reduce((count, project) => {
  if (!projectMatchesAllFiltersWith(project, draftCategoryFilter.value, searchQuery.value, rentFilter.value, availabilityFilter.value, sizeFilter.value, zoneFilter.value, readinessFilter.value, draftOperationalFilters())) return count
  return count + getProjectLayout(project).stalls.filter(stall => matchesStallFiltersWith(stall, rentFilter.value, availabilityFilter.value, sizeFilter.value, draftOperationalFilters())).length
}, 0))

const hasActiveFilters = computed(() => Boolean(
  appliedCategoryFilter.value || appliedSearchQuery.value.trim() || appliedRentFilter.value !== 'all' || appliedAvailabilityFilter.value !== 'all' || appliedSizeFilter.value !== 'all' || appliedZoneFilter.value !== 'all' || appliedReadinessFilter.value !== 'all' || hasOperationalFilter(appliedOperationalFilters())
))
const advancedFilterCount = computed(() => [
  appliedCategoryFilter.value,
  appliedRentFilter.value !== 'all',
  appliedSizeFilter.value !== 'all',
  !props.isOverviewPage && appliedAvailabilityFilter.value !== 'all',
  appliedZoneFilter.value !== 'all',
  appliedReadinessFilter.value !== 'all',
  appliedFloorFilter.value !== 'all',
  appliedContractExpiryFilter.value !== 'all',
  appliedBookingRangeFilter.value !== 'all',
  appliedRenewalStatusFilter.value !== 'all',
  appliedReturnStatusFilter.value !== 'all',
  appliedLegalStatusFilter.value !== 'all',
  appliedMaintenanceStatusFilter.value !== 'all',
  appliedRequestTypeFilter.value !== 'all'
].filter(Boolean).length)
const hasAdvancedFilters = computed(() => advancedFilterCount.value > 0)

// Category legend (read-only): สี → หมวดหมู่ → เลขโครงการ
const categoryStats = computed(() => {
  return categories.value
    .map(cat => {
      const projs = overviewProjects.value.filter(p => p.category === cat.name)
      return {
        ...cat,
        projectCount: projs.length,
        projectNumbers: projs.map(p => p.sectionNumber)
      }
    })
    .filter(cat => cat.projectCount > 0)
})

const categoryOptions = computed<SelectOption[]>(() => [
  { value: '', label: `ทุกหมวดหมู่ (${overviewProjects.value.length} โครงการ)` },
  ...categoryStats.value.map(cat => ({
    value: cat.name,
    label: `${cat.name} (${cat.projectCount})`
  }))
])

const zoneOptions = computed<SelectOption[]>(() => [
  { value: 'all', label: `ทุกโซน (${projects.value.length} โซน)` },
  ...projects.value.map(project => ({
    value: project.id,
    label: `${project.sectionNumber ? `โครงการ ${project.sectionNumber} · ` : ''}${project.name}`
  }))
])

const rentOptions: SelectOption[] = [
  { value: 'all', label: 'ทั้งหมด' },
  { value: 'under-2000', label: 'ต่ำกว่า ฿2,000' },
  { value: '2000-3000', label: '฿2,000–3,000' },
  { value: 'over-3000', label: 'มากกว่า ฿3,000' },
]

const sizeOptions: SelectOption[] = [
  { value: 'all', label: 'ทุกขนาด' },
  { value: 'up-to-5', label: 'ไม่เกิน 5 ตร.ม.' },
  { value: '5-to-10', label: 'มากกว่า 5–10 ตร.ม.' },
  { value: 'over-10', label: 'มากกว่า 10 ตร.ม.' }
]

const availabilityOptions: SelectOption[] = [
  { value: 'all', label: 'ทุกสถานะ' },
  { value: 'available', label: 'ว่าง' },
  { value: 'occupied', label: 'มีผู้เช่า' },
  { value: 'reserved', label: 'ติดจอง' },
]

const floorFilterOptions = computed<SelectOption[]>(() => [
  { value: 'all', label: 'ทุกชั้น' },
  ...marketFloors.value.map(floor => ({ value: String(floor.id), label: floor.name || `ชั้น ${floor.level}` }))
])

const contractStatusOptions: SelectOption[] = [
  { value: 'all', label: 'ทุกสถานะ' },
  { value: 'active', label: 'อยู่ในสัญญา' },
  { value: 'expiring', label: 'ใกล้หมดอายุ' },
  { value: 'expired', label: 'หมดอายุ' }
]
const contractStatusLabel = computed(() => contractStatusOptions.find(option => option.value === contractStatusFilter.value)?.label || '')
const contractExpiryOptions: SelectOption[] = [
  { value: 'all', label: 'ไม่จำกัด' },
  { value: '30', label: 'ภายใน 30 วัน' },
  { value: '60', label: 'ภายใน 60 วัน' },
  { value: '90', label: 'ภายใน 90 วัน' },
  { value: 'custom', label: 'กำหนดเอง (ช่วงวันที่)' }
]
const bookingRangeOptions: SelectOption[] = [
  { value: 'all', label: 'ไม่จำกัด' },
  { value: '7', label: '7 วันล่าสุด' },
  { value: '30', label: '30 วันล่าสุด' },
  { value: '90', label: '90 วันล่าสุด' },
  { value: 'custom', label: 'กำหนดเอง (ช่วงวันที่)' }
]
const renewalStatusOptions: SelectOption[] = [
  { value: 'all', label: 'ทุกสถานะ' },
  { value: 'awaiting', label: 'รอตอบรับ (รอต่อ + พิจารณา)' },
  { value: 'pending', label: 'รอต่ออายุ' },
  { value: 'reviewing', label: 'อยู่ระหว่างพิจารณา' },
  { value: 'approved', label: 'อนุมัติ' },
  { value: 'declined', label: 'ไม่ประสงค์ต่อ' }
]
const returnStatusOptions: SelectOption[] = [
  { value: 'all', label: 'ทุกสถานะ' },
  { value: 'notified', label: 'แจ้งไม่ต่อ' },
  { value: 'pending', label: 'รอส่งคืน' },
  { value: 'inspection', label: 'รอตรวจรับ' },
  { value: 'returned', label: 'คืนแล้ว' },
  { value: 'overdue', label: 'เกินกำหนดคืน' }
]
const legalStatusOptions: SelectOption[] = [
  { value: 'all', label: 'ทุกสถานะ' },
  { value: 'legal_review', label: 'รอส่งฝ่ายกฎหมาย' },
  { value: 'proceeding', label: 'อยู่ระหว่างคดี' },
  { value: 'enforcement', label: 'บังคับคดี' },
  { value: 'closed', label: 'สิ้นสุด' }
]
const maintenanceStatusFilterOptions: SelectOption[] = [
  { value: 'all', label: 'ทุกสถานะ' },
  { value: 'pending', label: 'รอซ่อม' },
  { value: 'repairing', label: 'กำลังซ่อม' },
  { value: 'completed', label: 'ซ่อมเสร็จ' },
  { value: 'overdue', label: 'เกินกำหนด' }
]
const requestTypeOptions: SelectOption[] = [
  { value: 'all', label: 'ทุกประเภท' },
  { value: 'contract_name_change', label: 'เปลี่ยนชื่อคู่สัญญา' },
  { value: 'transfer', label: 'โอนสิทธิ์' },
  { value: 'contract_cancel', label: 'ยกเลิกสัญญา' },
  { value: 'sublease', label: 'ผู้เช่าช่วง' }
]

const readinessOptions: SelectOption[] = [
  { value: 'all', label: 'ทั้งหมด' },
  { value: 'ready', label: 'พร้อมเผยแพร่' },
  { value: 'missing-stalls', label: 'ยังไม่มีแผงในโซน' },
  { value: 'incomplete-stalls', label: 'ข้อมูลแผงไม่ครบ' },
]

const handleCategoryChange = (val: string | number) => {
  draftCategoryFilter.value = val ? String(val) : ''
}

// ---------------- Project drill-down state ----------------
const hoveredProject = ref<ProjectZone | null>(null)
const hoveredStallKey = ref<string | null>(null)
const activeProjectId = ref<string | null>(null)
const isProjectPanelExiting = ref(false)
const overviewEditProjectId = ref<string | null>(null)
const isPolygonEditMode = ref(false)
const isStallPolygonEditMode = ref(false)
const editingStallCode = ref<string | null>(null)
const stallEditSnapshot = ref<ProjectZone | null>(null)
const stallEditHistory = ref<ProjectZone[]>([])
const isDrawingZone = ref(false)
const isDrawingStall = ref(false)
const isAddingFacility = ref(false)
const isDeletingStall = ref(false)
const stallDrawShape = ref<'rectangle' | 'polygon'>('rectangle')
const stallRectAlign = ref<'zone' | 'axis'>('zone')
const rectStartCorner = ref<[number, number] | null>(null)
const drawnStallDimensions = ref<{ widthM: number; depthM: number } | null>(null)
const drawingPoints = ref<[number, number][]>([])
const drawingCursor = ref<[number, number] | null>(null)
const marketFloorBeforeProject = ref('all')
const currentFloorId = ref<string>('1')
const floorSelectorProject = ref<ProjectZone | null>(null)

// จำโครงการที่แผงถูกเลือกไว้ (สำหรับปุ่มกลับมาดูแผงตัวเองใน Step 3)
const initialBookedProjectId = ref<string | null>(props.selectedProjectId || null)
watch(() => props.selectedProjectId, (newId) => {
  if (newId) initialBookedProjectId.value = newId
})

const bookedProjectId = computed(() => {
  if (props.selectedProjectId) return props.selectedProjectId
  if (props.selectedStallCodes && props.selectedStallCodes.length > 0) {
    for (const p of projects.value) {
      const allLayouts = getAllFloorLayoutsForProject(p)
      for (const layout of allLayouts.values()) {
        if (layout.stalls.some(s => props.selectedStallCodes.includes(s.code) || props.selectedStallCodes.includes(s.selectionId || ''))) {
          return p.id
        }
      }
    }
  }
  return initialBookedProjectId.value || null
})
const bookedProject = computed(() => {
  const id = bookedProjectId.value
  return (id ? projects.value.find(p => p.id === id) : null) || activeProject.value
})

const isNotViewingBookedProject = computed(() => {
  const targetId = bookedProjectId.value
  if (!targetId) return false
  return activeProjectId.value !== targetId
})

const activeProject = computed(() => projects.value.find(p => p.id === activeProjectId.value) || null)
const overviewEditProject = computed(() => {
  if (activeProject.value) return null
  const targetId = overviewEditProjectId.value
  return (targetId ? projects.value.find(p => p.id === targetId) : null) || null
})
const polygonEditTargetProject = computed<ProjectZone | null>(() => activeProject.value || overviewEditProject.value)
const polygonEditSnapshot = ref<ProjectZone[] | null>(null)
const polygonEditHistory = ref<ProjectZone[][]>([])

const rememberPolygonEdit = () => {
  polygonEditHistory.value.push(cloneProjects(currentMarket.value.projects))
}

const isAnyEditingToolActive = computed(() => (
  isPolygonEditMode.value ||
  isStallPolygonEditMode.value ||
  isDrawingZone.value ||
  isDrawingStall.value ||
  isAddingFacility.value ||
  isDeletingStall.value ||
  drawingPoints.value.length > 0
))

const togglePolygonEdit = () => {
  if (projects.value.length === 0) return
  pinnedFacilityPopup.value = false
  hoverPopup?.remove()
  isStallPolygonEditMode.value = false
  clearStallVertexHandles()
  isPolygonEditMode.value = !isPolygonEditMode.value
  if (isPolygonEditMode.value) {
    polygonEditSnapshot.value = cloneProjects(currentMarket.value.projects)
    polygonEditHistory.value = []
    if (!activeProject.value) {
      overviewEditProjectId.value = null
    }
    renderPolygonVertexHandles()
  } else {
    clearPolygonVertexHandles()
    overviewEditProjectId.value = null
    polygonEditSnapshot.value = null
    polygonEditHistory.value = []
  }
  refreshSources()
}

const savePolygonEdit = () => {
  const target = polygonEditTargetProject.value
  if (target) {
    const latest = currentMarket.value.projects.find(item => item.id === target.id)
    if (latest) emit('updateProjectPolygon', cloneProjects(latest))
  }
  isPolygonEditMode.value = false
  overviewEditProjectId.value = activeProjectId.value || null
  polygonEditSnapshot.value = null
  polygonEditHistory.value = []
  clearPolygonVertexHandles()
  refreshSources()
}

const undoPolygonEdit = () => {
  const previousProjects = polygonEditHistory.value.pop()
  if (!previousProjects) return
  currentMarket.value.projects = cloneProjects(previousProjects)
  clearLayoutCache()
  refreshSources()
  const target = polygonEditTargetProject.value
  if (target) {
    const restored = currentMarket.value.projects.find(item => item.id === target.id)
    if (restored) emit('updateProjectPolygon', cloneProjects(restored))
  }
  renderPolygonVertexHandles()
}

const cancelPolygonEdit = () => {
  const targetId = polygonEditTargetProject.value?.id
  const snapshot = polygonEditSnapshot.value
  if (snapshot) {
    currentMarket.value.projects = cloneProjects(snapshot)
    clearLayoutCache()
    if (targetId) {
      const restored = snapshot.find(item => item.id === targetId)
      if (restored) emit('updateProjectPolygon', cloneProjects(restored))
    }
  }
  isPolygonEditMode.value = false
  overviewEditProjectId.value = activeProjectId.value || null
  polygonEditSnapshot.value = null
  polygonEditHistory.value = []
  clearPolygonVertexHandles()
  refreshSources()
}

const toggleStallPolygonEdit = () => {
  if (!activeProject.value || !activeLayout.value?.stalls?.length) return
  pinnedFacilityPopup.value = false
  hoverPopup?.remove()
  isPolygonEditMode.value = false
  clearPolygonVertexHandles()
  isStallPolygonEditMode.value = !isStallPolygonEditMode.value
  if (isStallPolygonEditMode.value) {
    stallEditSnapshot.value = cloneProjects(currentMarket.value.projects.find(item => item.id === activeProject.value?.id) || null)
    stallEditHistory.value = []
    if (!editingStallCode.value || !activeLayout.value.stalls.some(s => s.code === editingStallCode.value)) {
      editingStallCode.value = activeLayout.value.stalls[0]?.code || null
    }
    renderStallVertexHandles()
  } else {
    clearStallVertexHandles()
    stallEditSnapshot.value = null
    stallEditHistory.value = []
  }
  refreshSources()
}

const saveStallPolygonEdit = () => {
  const project = currentMarket.value.projects.find(item => item.id === activeProject.value?.id)
  const stall = project?.stalls?.find((item: any) => String(item.code) === String(editingStallCode.value))
  if (project && stall) emit('updateStallPolygon', cloneProjects(project), cloneProjects(stall))
  toggleStallPolygonEdit()
}

const rememberStallEdit = () => {
  const project = currentMarket.value.projects.find(item => item.id === activeProject.value?.id)
  if (!project) return
  stallEditHistory.value.push(cloneProjects(project))
}

const emitCurrentEditingStallUpdate = () => {
  const project = currentMarket.value.projects.find(item => item.id === activeProject.value?.id)
  const stall = project?.stalls?.find((item: any) => String(item.code) === String(editingStallCode.value))
  if (project && stall) emit('updateStallPolygon', cloneProjects(project), cloneProjects(stall))
}

const undoStallPolygonEdit = () => {
  const previous = stallEditHistory.value.pop()
  if (!previous) return
  const index = currentMarket.value.projects.findIndex(item => item.id === previous.id)
  if (index >= 0) currentMarket.value.projects[index] = cloneProjects(previous)
  clearLayoutCache()
  refreshSources()
  emitCurrentEditingStallUpdate()
  renderStallVertexHandles()
}

const cancelStallPolygonEdit = () => {
  const snapshot = stallEditSnapshot.value
  if (snapshot) {
    const index = currentMarket.value.projects.findIndex(item => item.id === snapshot.id)
    if (index >= 0) currentMarket.value.projects[index] = cloneProjects(snapshot)
    clearLayoutCache()
  }
  emitCurrentEditingStallUpdate()
  isStallPolygonEditMode.value = false
  stallEditSnapshot.value = null
  stallEditHistory.value = []
  clearStallVertexHandles()
  refreshSources()
}

const requestDeleteActiveProject = () => {
  const target = activeProject.value || overviewEditProject.value
  if (!target) return
  const layoutSummary = getProjectLayout(target)
  const stats = layoutSummary
    ? {
        total: layoutSummary.total,
        available: layoutSummary.available,
        reserved: layoutSummary.reserved,
        occupied: layoutSummary.occupied
      }
    : undefined
  emit('requestDeleteProject', target, stats)
}

const getZoneReferenceAngleRad = (p1: [number, number]): number => {
  if (stallRectAlign.value === 'axis' || !activeProject.value?.coordinates?.length) return 0
  const ring = activeProject.value.coordinates
  if (ring.length < 2) return 0
  const latRad = (p1[1] * Math.PI) / 180
  const mPerDegLng = 111320 * Math.cos(latRad)
  const mPerDegLat = 110540

  let bestDistSq = Infinity
  let bestAngle = 0
  for (let i = 0; i < ring.length - 1; i++) {
    const a = ring[i]!
    const b = ring[i + 1]!
    const ax = (a[0] - p1[0]) * mPerDegLng
    const ay = (a[1] - p1[1]) * mPerDegLat
    const bx = (b[0] - p1[0]) * mPerDegLng
    const by = (b[1] - p1[1]) * mPerDegLat
    const vx = bx - ax
    const vy = by - ay
    const lenSq = vx * vx + vy * vy
    if (lenSq < 1e-4) continue
    const t = Math.max(0, Math.min(1, -(ax * vx + ay * vy) / lenSq))
    const px = ax + t * vx
    const py = ay + t * vy
    const distSq = px * px + py * py
    if (distSq < bestDistSq) {
      bestDistSq = distSq
      bestAngle = Math.atan2(vy, vx)
    }
  }
  return bestAngle
}

const buildRectangleCorners = (p1: [number, number], p2: [number, number]): { corners: [number, number][]; widthM: number; depthM: number } => {
  const latRad = (p1[1] * Math.PI) / 180
  const mPerDegLng = 111320 * Math.cos(latRad)
  const mPerDegLat = 110540
  const dx = (p2[0] - p1[0]) * mPerDegLng
  const dy = (p2[1] - p1[1]) * mPerDegLat
  const theta = getZoneReferenceAngleRad(p1)
  const ux = Math.cos(theta)
  const uy = Math.sin(theta)
  const vx = -uy
  const vy = ux
  const uDist = dx * ux + dy * uy
  const vDist = dx * vx + dy * vy

  const toLngLat = (mx: number, my: number): [number, number] => [
    p1[0] + mx / mPerDegLng,
    p1[1] + my / mPerDegLat
  ]

  const corners: [number, number][] = [
    [p1[0], p1[1]],
    toLngLat(uDist * ux, uDist * uy),
    toLngLat(dx, dy),
    toLngLat(vDist * vx, vDist * vy)
  ]

  return {
    corners,
    widthM: Math.max(0.5, Number(Math.abs(uDist).toFixed(1))),
    depthM: Math.max(0.5, Number(Math.abs(vDist).toFixed(1)))
  }
}

const liveRectPreview = computed<{ corners: [number, number][]; widthM: number; depthM: number } | null>(() => {
  if (!isDrawingStall.value || stallDrawShape.value !== 'rectangle') return null
  if (!rectStartCorner.value || !drawingCursor.value || drawingPoints.value.length >= 4) return null
  return buildRectangleCorners(rectStartCorner.value, drawingCursor.value)
})

const getDrawingZoneGeoJSON = () => {
  if (isDrawingStall.value && stallDrawShape.value === 'rectangle') {
    if (drawingPoints.value.length === 4) {
      const polygon = [...drawingPoints.value, drawingPoints.value[0]!]
      const drawnPointFeatures = drawingPoints.value.map(point => ({ type: 'Feature', geometry: { type: 'Point', coordinates: point }, properties: {} }))
      return {
        type: 'FeatureCollection',
        features: [
          { type: 'Feature', geometry: { type: 'Polygon', coordinates: [polygon] }, properties: {} },
          ...drawnPointFeatures
        ]
      }
    }
    if (rectStartCorner.value && drawingCursor.value) {
      const { corners } = buildRectangleCorners(rectStartCorner.value, drawingCursor.value)
      const polygon = [...corners, corners[0]!]
      const drawnPointFeatures = corners.map(point => ({ type: 'Feature', geometry: { type: 'Point', coordinates: point }, properties: {} }))
      return {
        type: 'FeatureCollection',
        features: [
          { type: 'Feature', geometry: { type: 'Polygon', coordinates: [polygon] }, properties: {} },
          ...drawnPointFeatures
        ]
      }
    }
    if (rectStartCorner.value) {
      return {
        type: 'FeatureCollection',
        features: [{ type: 'Feature', geometry: { type: 'Point', coordinates: rectStartCorner.value }, properties: {} }]
      }
    }
    return { type: 'FeatureCollection', features: [] }
  }

  const points = [...drawingPoints.value]
  if (drawingCursor.value && points.length) points.push(drawingCursor.value)
  const drawnPointFeatures = drawingPoints.value.map(point => ({ type: 'Feature', geometry: { type: 'Point', coordinates: point }, properties: {} }))
  if (points.length < 2) return { type: 'FeatureCollection', features: drawnPointFeatures }
  const polygon = points.length >= 3 ? [...points, points[0]!] : points
  return {
    type: 'FeatureCollection',
    features: [
      { type: 'Feature', geometry: { type: points.length >= 3 ? 'Polygon' : 'LineString', coordinates: points.length >= 3 ? [polygon] : polygon }, properties: {} },
      ...drawnPointFeatures
    ]
  }
}

const refreshDrawingZone = () => {
  map?.getSource('zone-drawing')?.setData(getDrawingZoneGeoJSON())
}

const setStallDrawShape = (mode: 'rectangle' | 'polygon') => {
  if (stallDrawShape.value === mode) return
  stallDrawShape.value = mode
  rectStartCorner.value = null
  drawnStallDimensions.value = null
  drawingPoints.value = []
  refreshDrawingZone()
}

const toggleStallRectAlign = () => {
  stallRectAlign.value = stallRectAlign.value === 'zone' ? 'axis' : 'zone'
  if (rectStartCorner.value && drawingPoints.value.length === 4 && drawingPoints.value[2]) {
    const { corners, widthM, depthM } = buildRectangleCorners(rectStartCorner.value, drawingPoints.value[2])
    drawingPoints.value = corners
    drawnStallDimensions.value = { widthM, depthM }
  }
  refreshDrawingZone()
}

const beginDrawingZone = () => {
  pinnedFacilityPopup.value = false
  hoverPopup?.remove()
  isDrawingZone.value = true
  isDrawingStall.value = false
  isAddingFacility.value = false
  isPolygonEditMode.value = false
  isStallPolygonEditMode.value = false
  clearPolygonVertexHandles()
  clearStallVertexHandles()
  rectStartCorner.value = null
  drawnStallDimensions.value = null
  drawingPoints.value = []
  drawingCursor.value = null
  refreshDrawingZone()
  if (map) map.getCanvas().style.cursor = 'crosshair'
}

const beginDrawingStall = () => {
  if (!activeProject.value) return
  pinnedFacilityPopup.value = false
  hoverPopup?.remove()
  isDrawingStall.value = true
  isDrawingZone.value = false
  isAddingFacility.value = false
  isPolygonEditMode.value = false
  isStallPolygonEditMode.value = false
  clearPolygonVertexHandles()
  clearStallVertexHandles()
  rectStartCorner.value = null
  drawnStallDimensions.value = null
  drawingPoints.value = []
  drawingCursor.value = null
  refreshDrawingZone()
  if (map) map.getCanvas().style.cursor = 'crosshair'
}

const beginAddingFacility = () => {
  pinnedFacilityPopup.value = false
  hoverPopup?.remove()
  isAddingFacility.value = true
  isDrawingZone.value = false
  isDrawingStall.value = false
  isPolygonEditMode.value = false
  isStallPolygonEditMode.value = false
  clearPolygonVertexHandles()
  clearStallVertexHandles()
  if (map) map.getCanvas().style.cursor = 'crosshair'
}

const undoDrawingPoint = () => {
  if (isDrawingStall.value && stallDrawShape.value === 'rectangle') {
    rectStartCorner.value = null
    drawnStallDimensions.value = null
    drawingPoints.value = []
    refreshDrawingZone()
    return
  }
  drawingPoints.value.pop()
  refreshDrawingZone()
}

const materializeActiveProjectStalls = (): ProjectZone | null => {
  if (!activeProject.value) return null
  const projectID = activeProject.value.id
  const sourceProject = currentMarket.value.projects.find(item => item.id === projectID)
  if (!sourceProject) return null
  if (!sourceProject.stalls || sourceProject.stalls.length === 0) {
    const allGenerated: any[] = []
    const floors = sourceProject.floors?.length ? sourceProject.floors : [{ id: '1', name: 'ชั้น 1', shortLabel: '1F', level: 1 }]
    floors.forEach(fl => {
      const flLayout = applyMarketStatus(getProjectStallLayout(sourceProject, fl.id))
      flLayout.stalls.forEach(st => {
        allGenerated.push({ ...st, floor: String(fl.id) })
      })
    })
    sourceProject.stalls = cloneProjects(allGenerated)
    sourceProject.layoutMode = 'manual'
    clearLayoutCache()
  }
  return sourceProject
}

const triggerDeleteStall = (stall: any) => {
  const sourceProject = materializeActiveProjectStalls()
  if (!sourceProject || !stall) return
  isDeletingStall.value = false
  if (isStallPolygonEditMode.value) {
    isStallPolygonEditMode.value = false
    clearStallVertexHandles()
  }
  if (map) map.getCanvas().style.cursor = ''
  emit('requestDeleteStall', cloneProjects(sourceProject), cloneProjects(stall))
}

const beginDeletingStall = () => {
  if (!activeProject.value || !activeLayout.value?.stalls?.length) return
  pinnedFacilityPopup.value = false
  hoverPopup?.remove()
  isDrawingZone.value = false
  isDrawingStall.value = false
  isAddingFacility.value = false
  isPolygonEditMode.value = false
  isStallPolygonEditMode.value = false
  clearPolygonVertexHandles()
  clearStallVertexHandles()
  if (activeLayout.value.stalls.length === 1) {
    triggerDeleteStall(activeLayout.value.stalls[0])
    return
  }
  isDeletingStall.value = true
  if (map) map.getCanvas().style.cursor = 'pointer'
}

const deleteCurrentlyEditingStall = () => {
  if (!editingStallCode.value || !activeLayout.value?.stalls?.length) return
  const target = activeLayout.value.stalls.find(s => String(s.code) === String(editingStallCode.value))
  if (target) triggerDeleteStall(target)
}

const cancelDrawing = () => {
  isDrawingZone.value = false
  isDrawingStall.value = false
  isAddingFacility.value = false
  isDeletingStall.value = false
  rectStartCorner.value = null
  drawnStallDimensions.value = null
  drawingPoints.value = []
  drawingCursor.value = null
  refreshDrawingZone()
  if (map) map.getCanvas().style.cursor = ''
}

const finishDrawing = () => {
  if (drawingPoints.value.length < 3) return
  const points = drawingPoints.value.map(point => [...point] as [number, number])
  const project = activeProject.value
  const isStall = isDrawingStall.value
  const dims = drawnStallDimensions.value ? { ...drawnStallDimensions.value } : undefined
  cancelDrawing()
  if (isStall && project) emit('requestAddStallLayout', project, points, currentFloorId.value, dims)
  else emit('requestAddProject', points)
}

// Multi-floor management for active project
const activeProjectFloors = computed<ProjectFloor[]>(() => {
  if (!activeProject.value?.floors || activeProject.value.floors.length <= 1) return []
  return activeProject.value.floors
})

const activeProjectFloorsReversed = computed<ProjectFloor[]>(() => {
  // Physical building stack: highest floor on top (3F, 2F, 1F)
  return [...activeProjectFloors.value].sort((a, b) => b.level - a.level)
})

const currentFloorObject = computed<ProjectFloor | null>(() => {
  if (!activeProject.value?.floors || activeProject.value.floors.length === 0) return null
  return activeProject.value.floors.find(f => f.id === currentFloorId.value) || activeProject.value.floors[0] || null
})

const activeLayout = computed(() => {
  if (!activeProjectId.value || !activeProject.value) return null
  return applyMarketStatus(getProjectStallLayout(activeProject.value, currentFloorId.value))
})

const getFloorAvailCount = (floorId: string): number => {
  if (!activeProject.value) return 0
  const layout = applyMarketStatus(getProjectStallLayout(activeProject.value, floorId))
  return layout.available
}

const getFloorReservedCount = (floorId: string): number => {
  if (!activeProject.value) return 0
  const layout = applyMarketStatus(getProjectStallLayout(activeProject.value, floorId))
  return layout.reserved
}

const suggestedAvailableFloor = computed<ProjectFloor | null>(() => {
  if (!activeProject.value?.floors || activeProject.value.floors.length <= 1) return null
  return activeProject.value.floors.find(f => f.id !== currentFloorId.value && getFloorAvailCount(f.id) > 0) || null
})

const getFloorLayout = (project: ProjectZone, floorId: string) => {
  return applyMarketStatus(getProjectStallLayout(project, floorId))
}

const switchFloor = (floorId: string) => {
  if (currentFloorId.value === floorId) return
  currentFloorId.value = floorId
  pinnedFacilityPopup.value = false
  setHoveredStall(null)
  hoverPopup?.remove()
  refreshSources()
  renderSelectedStallPins()
  if (isStallPolygonEditMode.value) {
    editingStallCode.value = activeLayout.value?.stalls[0]?.code || null
    renderStallVertexHandles()
  }
}

const currentProjectInfo = computed<ProjectZone | null>(() => {
  return activeProject.value || hoveredProject.value || overviewProjects.value[0] || null
})

const currentProjectLayout = computed(() => (currentProjectInfo.value ? getProjectLayout(currentProjectInfo.value) : null))

// Stall search and filter states in side panel
const stallSearchQuery = ref('')
const stallFilterTab = ref<'all' | 'available' | 'reserved'>('all')

const filteredStalls = computed(() => {
  if (!activeLayout.value?.stalls) return []
  let list = activeLayout.value.stalls
  if (stallFilterTab.value === 'available') {
    list = list.filter(s => s.status === 'available')
  } else if (stallFilterTab.value === 'reserved') {
    list = list.filter(s => s.status === 'reserved')
  }
  if (stallSearchQuery.value.trim()) {
    const q = stallSearchQuery.value.trim().toLowerCase()
    list = list.filter(s => s.code.toLowerCase().includes(q) || String(s.name || '').toLowerCase().includes(q))
  }
  return list
})

const getStallSelectionKey = (stall: MapStall): string => stall.selectionId || stall.code
const getStallMapLabel = (stall: MapStall): string => {
  const sourceNumber = stall.sourceProperties?.number ?? stall.sourceProperties?.Number
  return String(sourceNumber ?? '').trim()
}
const getStallDisplayName = (stall: MapStall): string => getStallMapLabel(stall) || String(stall.name || '').trim()
const isStallSelected = (stall: MapStall): boolean => {
  const selectedKeys = props.selectedStallKeys?.length ? props.selectedStallKeys : props.selectedStallCodes
  const normalizedSelectedKeys = new Set((selectedKeys || []).map(key => String(key).trim().toUpperCase()))
  return [stall.selectionId, stall.code]
    .filter(Boolean)
    .some(key => normalizedSelectedKeys.has(String(key).trim().toUpperCase()))
}

const setHoveredStall = (stall: MapStall | null): boolean => {
  const nextKey = stall ? getStallSelectionKey(stall) : null
  if (hoveredStallKey.value === nextKey) return false

  const layout = activeLayout.value
  const previousKey = hoveredStallKey.value
  if (map?.getSource('project-stalls') && layout && previousKey) {
    const previousIndex = layout.stalls.findIndex(item => getStallSelectionKey(item) === previousKey)
    if (previousIndex >= 0) {
      map.setFeatureState({ source: 'project-stalls', id: previousIndex + 1 }, { isHovered: false })
    }
  }
  hoveredStallKey.value = nextKey
  if (map?.getSource('project-stalls') && layout && nextKey) {
    const nextIndex = layout.stalls.findIndex(item => getStallSelectionKey(item) === nextKey)
    if (nextIndex >= 0) {
      map.setFeatureState({ source: 'project-stalls', id: nextIndex + 1 }, { isHovered: true })
    }
  }
  return true
}

const handleStallCardHover = (stall: MapStall) => {
  if (!map || !hoverPopup) return
  setHoveredStall(stall)
  const isMyBooked = !props.previewMode && (props.myBookedStallCodes || []).includes(stall.code)
  const canClick = props.managementMode || (stall.status !== 'occupied' && !isMyBooked)
  map.getCanvas().style.cursor = canClick ? 'pointer' : 'not-allowed'
  showPopupForStall(stall)
}

const handleStallCardLeave = () => {
  if (!map) return
  setHoveredStall(null)
  map.getCanvas().style.cursor = ''
  hoverPopup?.remove()
}

const handleStallCardClick = (stall: MapStall) => {
  if (stall.status === 'occupied' && !props.managementMode) return
  emit('selectStall', stall)
}

watch([stallSearchQuery, stallFilterTab], () => {
  refreshSources()
})

// ================================================================
// MAPLIBRE GL JS ENGINE IMPLEMENTATION
// ================================================================
let maplibregl: any = null
let map: any = null
let projectMarkers: any[] = []
let polygonHandleMarkers: any[] = []
let stallHandleMarkers: any[] = []
let facilityMarkers: any[] = []
const pinnedFacilityPopup = ref(false)
const projectMarkerElements = new Map<string, HTMLElement>()
let hoverPopup: any = null

// Interior label point (inside the polygon frame, not the vertex average)
const projectToLngLat = (p: ProjectZone): [number, number] => projectLayoutCache.get(p.id)?.labelPoint || p.centroid

const getProjectBounds = (p: ProjectZone): [[number, number], [number, number]] => {
  const lngs = p.coordinates.map(c => c[0])
  const lats = p.coordinates.map(c => c[1])
  return [[Math.min(...lngs), Math.min(...lats)], [Math.max(...lngs), Math.max(...lats)]]
}

const getStallsBounds = (stalls: MapStall[]): [[number, number], [number, number]] | null => {
  if (stalls.length === 0) return null
  const allCoords = stalls.flatMap(s => s.coordinates)
  const lngs = allCoords.map(c => c[0])
  const lats = allCoords.map(c => c[1])
  const minLng = Math.min(...lngs)
  const minLat = Math.min(...lats)
  const maxLng = Math.max(...lngs)
  const maxLat = Math.max(...lats)
  const pad = 0.00008
  return [[minLng - pad, minLat - pad], [maxLng + pad, maxLat + pad]]
}

// 3 Real Tile Styles (Watermark-free, zero API key required)
const standardStyle: any = {
  version: 8,
  sources: {
    'osm-street': {
      type: 'raster',
      tiles: [
        'https://tile.openstreetmap.org/{z}/{x}/{y}.png'
      ],
      tileSize: 256,
      attribution: '&copy; OpenStreetMap contributors',
      maxzoom: 19
    }
  },
  layers: [
    {
      id: 'osm-street-layer',
      type: 'raster',
      source: 'osm-street',
      minzoom: 0,
      maxzoom: 22
    }
  ]
}

const darkStyle: any = {
  version: 8,
  sources: {
    // ใช้ไทล์ OpenStreetMap ชุดเดียวกับแผนที่ปกติ (ตัวหนังสือภาษาไทยคมและซูมได้ถึงระดับ 19) แล้วกลับความสว่างด้วย raster paint เป็นโทนมืด
    // เลิกใช้ไทล์ Esri/CARTO ของเดิม: Esri จำกัดซูม 16 ทำให้ตัวหนังสือเบลอ ส่วน CARTO ต้องมี API key
    'osm-dark': {
      type: 'raster',
      tiles: [
        'https://tile.openstreetmap.org/{z}/{x}/{y}.png'
      ],
      tileSize: 256,
      attribution: '&copy; OpenStreetMap contributors',
      maxzoom: 19
    }
  },
  layers: [
    {
      id: 'osm-dark-layer',
      type: 'raster',
      source: 'osm-dark',
      minzoom: 0,
      maxzoom: 22,
      paint: {
        'raster-brightness-min': 0.9,
        'raster-brightness-max': 0.1,
        'raster-hue-rotate': 180,
        'raster-saturation': 0.15,
        'raster-contrast': 0.15
      }
    }
  ]
}

const satelliteStyle: any = {
  version: 8,
  sources: {
    'esri-satellite': {
      type: 'raster',
      tiles: [
        'https://server.arcgisonline.com/ArcGIS/rest/services/World_Imagery/MapServer/tile/{z}/{y}/{x}'
      ],
      tileSize: 256,
      attribution: 'Tiles &copy; Esri',
      maxzoom: 19
    },
    'esri-labels': {
      type: 'raster',
      tiles: [
        'https://server.arcgisonline.com/ArcGIS/rest/services/Reference/World_Boundaries_and_Places/MapServer/tile/{z}/{y}/{x}'
      ],
      tileSize: 256,
      attribution: 'Tiles &copy; Esri',
      maxzoom: 19
    }
  },
  layers: [
    {
      id: 'esri-satellite-layer',
      type: 'raster',
      source: 'esri-satellite',
      minzoom: 0,
      maxzoom: 22
    },
    {
      id: 'esri-labels-layer',
      type: 'raster',
      source: 'esri-labels',
      minzoom: 0,
      maxzoom: 22
    }
  ]
}

const getStyleConfig = (styleName?: string) => {
  if (styleName === 'satellite') return satelliteStyle
  if (styleName === 'dark') return darkStyle
  return standardStyle
}

// ---------------- 3D Buildings (MapLibre "Display buildings in 3D") ----------------
// ใช้ vector tiles ของ OpenFreeMap (OpenMapTiles schema, ไม่ต้องใช้ API key)
const BUILDINGS_SOURCE_ID = 'openfreemap'
const BUILDINGS_LAYER_ID = '3d-buildings'
const DEFAULT_PITCH = 0
const DEFAULT_BEARING = 0
const DRILL_PITCH = 0
const BUILDINGS_OPACITY = 0.85
const BUILDINGS_OPACITY_DRILL = 0.25

const getBuildingColorExpr = (styleName?: string): any => {
  if (styleName === 'dark') {
    return ['interpolate', ['linear'], ['get', 'render_height'], 0, '#334155', 60, '#475569', 200, '#64748b']
  }
  if (styleName === 'satellite') {
    return ['interpolate', ['linear'], ['get', 'render_height'], 0, '#e2e8f0', 60, '#cbd5e1', 200, '#94a3b8']
  }
  return ['interpolate', ['linear'], ['get', 'render_height'], 0, '#d6dce5', 40, '#a5b4fc', 120, '#6366f1', 300, '#93c5fd']
}

const add3DBuildings = () => {
  if (!map) return
  try {
    if (map.getLayer(BUILDINGS_LAYER_ID)) map.removeLayer(BUILDINGS_LAYER_ID)
    if (currentStyle.value === 'satellite') return
    if (!map.getSource(BUILDINGS_SOURCE_ID)) {
      map.addSource(BUILDINGS_SOURCE_ID, {
        type: 'vector',
        url: 'https://tiles.openfreemap.org/planet',
        attribution: '&copy; OpenFreeMap &copy; OpenMapTiles'
      })
    }
    // วางอาคาร 3D ไว้ใต้กรอบโครงการ
    const beforeId = map.getLayer('projects-fill') ? 'projects-fill' : undefined
    map.addLayer({
      id: BUILDINGS_LAYER_ID,
      source: BUILDINGS_SOURCE_ID,
      'source-layer': 'building',
      type: 'fill-extrusion',
      minzoom: 15,
      filter: ['!=', ['get', 'hide_3d'], true],
      paint: {
        'fill-extrusion-color': getBuildingColorExpr(resolveStyle(currentStyle.value)),
        'fill-extrusion-height': [
          'interpolate', ['linear'], ['zoom'],
          15, 0,
          16, ['coalesce', ['get', 'render_height'], 4]
        ],
        'fill-extrusion-base': [
          'interpolate', ['linear'], ['zoom'],
          15, 0,
          16, ['coalesce', ['get', 'render_min_height'], 0]
        ],
        'fill-extrusion-opacity': activeProjectId.value ? BUILDINGS_OPACITY_DRILL : BUILDINGS_OPACITY
      }
    }, beforeId)
  } catch (err) {
    console.warn('[map] 3D buildings layer failed:', err)
  }
}

const setBuildingsOpacity = (value: number) => {
  if (map?.getLayer(BUILDINGS_LAYER_ID)) {
    map.setPaintProperty(BUILDINGS_LAYER_ID, 'fill-extrusion-opacity', value)
  }
}

const mapContainerRef = ref<HTMLElement | null>(null)
const isMapLoading = ref(true)
// The registry starts every market with no projects; polygons only exist once /gis-config
// has been fetched and drawn. Keep the skeleton up until both the basemap and that data are ready.
let isMapBaseReady = false
let isGisHydrating = true
const revealMapIfReady = () => {
  if (isMapBaseReady && !isGisHydrating) isMapLoading.value = false
}

const initMap = () => {
  const container = (mapContainerRef.value || document.getElementById(props.containerId)) as any
  if (!container || map || !maplibregl) return

  isMapLoading.value = true

  if (container.__mapInstance) {
    try { container.__mapInstance.remove() } catch {}
    container.__mapInstance = null
  }
  container.innerHTML = ''

  map = new maplibregl.Map({
    container,
    style: getStyleConfig(resolveStyle(currentStyle.value)),
    center: currentMarketCenter.value,
    zoom: props.previewMode ? (currentMarketZoom.value + 1.2) : currentMarketZoom.value,
    minZoom: 13.0,
    maxZoom: 21,
    maxPitch: 65,
    pitch: currentMarketPitch.value,
    bearing: currentMarketBearing.value,
    canvasContextAttributes: { antialias: true },
    attributionControl: { compact: true } // ย่อเครดิตเหลือปุ่ม i ไม่ให้แถบยาวบังปุ่มมุมล่างบนจอเล็ก
  })
  container.__mapInstance = map

  // Standard Navigation Controls (Zoom, Compass)
  // ชุดปุ่มควบคุมมาตรฐานของระบบ (ตำแหน่งของฉัน / ซูม / เข็มทิศ / ปรับมุมมอง) ใช้ร่วมกันทุกแผนที่
  const navGroup = addStandardMapControls(maplibregl, map, {
    onFit: () => recenterMapView(),
    fitTitle: 'ปรับมุมมองให้เห็นทั้งพื้นที่',
    notify: showMapWarning
  })
  if (navGroup) watchCtrlBaseline(navGroup)

  hoverPopup = new maplibregl.Popup({
    closeButton: false,
    closeOnClick: false,
    offset: 8,
    maxWidth: '380px',
    className: 'mlp-popup'
  })

  bindMapEvents()

  let layersInitialized = false
  const initLayersOnce = () => {
    if (layersInitialized || !map) return
    layersInitialized = true
    setupProjectLayers()
    renderFacilityMarkers()
    const targetProject = bookedProject.value || (props.selectedProjectId ? projects.value.find(p => p.id === props.selectedProjectId) : null)
    if (targetProject) {
      let targetFloorId: string | undefined
      if (targetProject.floors && targetProject.floors.length > 0 && props.selectedStallCodes?.length) {
        for (const fl of targetProject.floors) {
          const flLayout = applyMarketStatus(getProjectStallLayout(targetProject, fl.id))
          if (flLayout.stalls.some(isStallSelected)) {
            targetFloorId = fl.id
            break
          }
        }
      }
      enterProject(targetProject, targetFloorId)
    }
    renderSelectedStallPins()

    // Wait for the initial viewport to settle before revealing the live map.
    // This keeps the management canvas from flashing as an empty white area
    // while MapLibre is still fetching its style and tiles.
    const clearLoading = () => {
      isMapBaseReady = true
      revealMapIfReady()
    }
    map.once('idle', clearLoading)
    setTimeout(clearLoading, 450)
  }

  if (map.isStyleLoaded()) {
    initLayersOnce()
  } else {
    map.once('style.load', initLayersOnce)
    map.once('load', initLayersOnce)
  }
}

// ---------------- GeoJSON builders ----------------
const closePolygonRing = (ring: [number, number][]): [number, number][] => {
  if (!ring.length) return []
  const first = ring[0]!
  const last = ring[ring.length - 1]!
  if (first[0] === last[0] && first[1] === last[1]) return ring
  return [...ring, [first[0], first[1]]]
}

const isProjectDimmed = (p: ProjectZone) => {
  // A project is only prominent when it has at least one stall that matches
  // every active filter. Previously this considered only category/search, so
  // filtering by status or rent left unrelated projects fully saturated.
  return !projectMatchesAllFilters(p)
}

const getProjectsGeoJson = () => {
  return {
    type: 'FeatureCollection',
    features: overviewProjects.value.map((p, index) => {
      const isOverviewEditTarget = isPolygonEditMode.value && !activeProjectId.value && overviewEditProjectId.value === p.id
      const isActive = activeProjectId.value === p.id || isOverviewEditTarget
      const isOtherWhileActive = Boolean(activeProjectId.value) && !isActive
      const isDimmedByFilter = isProjectDimmed(p)

      const fillColor = p.color
      const strokeColor = p.color

      return {
        type: 'Feature',
        id: index + 1,
        geometry: {
          type: 'Polygon',
          coordinates: [closePolygonRing(p.coordinates)]
        },
        properties: {
          id: p.id,
          fillColor,
          strokeColor: (isActive && !isDimmedByFilter) || (!isDimmedByFilter && hoveredProject.value?.id === p.id) ? '#2563eb' : strokeColor,
          isActive: isActive && !isDimmedByFilter,
          isHovered: !isDimmedByFilter && hoveredProject.value?.id === p.id,
          isDimmed: isDimmedByFilter || isOtherWhileActive
        }
      }
    })
  }
}

const getStallsGeoJson = () => {
  const layout = activeLayout.value
  const hasStallSearchMatches = Boolean(appliedSearchQuery.value.trim())
    && Boolean(layout?.stalls.some(stall => stallMatchesSearchWith(stall, appliedSearchQuery.value)))
  return {
    type: 'FeatureCollection',
    features: (layout?.stalls || []).map((s, idx) => {
      const isSearchMatch = Boolean(appliedSearchQuery.value.trim()) && stallMatchesSearchWith(s, appliedSearchQuery.value)
      const isSelected = isStallSelected(s)
      return {
        type: 'Feature',
        id: idx + 1,
        geometry: {
          type: 'Polygon',
          coordinates: [closePolygonRing(s.coordinates), ...(s.polygonHoles || []).map(closePolygonRing)]
        },
        properties: {
          // Survey exports can reuse a human-readable stall code. Keep the
          // layout index so map hit-testing always resolves the exact polygon.
          stallIndex: idx,
          code: s.code,
          label: getStallMapLabel(s),
          status: s.status,
          isSelected,
          isHovered: hoveredStallKey.value === getStallSelectionKey(s),
          floor: s.floor,
          floorName: s.floorName,
          isSearchMatch,
          isDimmed: !matchesStallFilters(s) || (hasStallSearchMatches && !stallMatchesSearchWith(s, appliedSearchQuery.value))
        }
      }
    })
  }
}

// ---------------- Layers ----------------
const setupProjectLayers = () => {
  if (!map) return

  ;['zone-drawing-points', 'zone-drawing-line', 'zone-drawing-fill', 'facilities-label', 'facilities-circle', 'facilities-polygon-line', 'facilities-polygon-fill', 'projects-hover-outline', 'projects-hover-casing', 'projects-active-outline', 'stalls-selected-stroke', 'stalls-selected-casing', 'stalls-hover-stroke', 'stalls-hover-casing', 'stalls-search-stroke', 'stalls-search-casing', 'stalls-label', 'stalls-line', 'stalls-casing', 'stalls-fill', 'projects-stroke', 'projects-casing', 'projects-fill'].forEach(id => {
    if (map.getLayer(id)) map.removeLayer(id)
  })
  ;['projects-polygons', 'project-stalls', 'market-facilities', 'zone-drawing'].forEach(id => {
    if (map.getSource(id)) map.removeSource(id)
  })

  map.addSource('projects-polygons', { type: 'geojson', data: getProjectsGeoJson() })
  map.addSource('project-stalls', { type: 'geojson', data: getStallsGeoJson() })
  map.addSource('market-facilities', { type: 'geojson', data: getFacilitiesGeoJson() })
  map.addSource('zone-drawing', { type: 'geojson', data: getDrawingZoneGeoJSON() })

  // กรอบพื้นที่โครงการ (Fill)
  map.addLayer({
    id: 'projects-fill',
    type: 'fill',
    source: 'projects-polygons',
    paint: {
      'fill-color': ['get', 'fillColor'],
      'fill-opacity': [
        'case',
        ['boolean', ['get', 'isActive'], false], 0.14,
        ['boolean', ['get', 'isHovered'], false], 0.24,
        ['boolean', ['get', 'isDimmed'], false], 0.1,
        0.42
      ]
    }
  })

  const facilityColor = [
    'match', ['get', 'type'],
    'sanitation', '#0284c7',
    'access', '#d97706',
    'common-area', '#0f766e',
    'building', '#64748b',
    'service', '#7c3aed',
    '#2563eb'
  ] as any
  map.addLayer({
    id: 'facilities-polygon-fill',
    type: 'fill',
    source: 'market-facilities',
    filter: ['==', ['geometry-type'], 'Polygon'],
    paint: { 'fill-color': facilityColor, 'fill-opacity': 0.34 }
  })
  map.addLayer({
    id: 'facilities-polygon-line',
    type: 'line',
    source: 'market-facilities',
    filter: ['==', ['geometry-type'], 'Polygon'],
    paint: { 'line-color': facilityColor, 'line-width': 2.2, 'line-opacity': 0.95 }
  })
  map.addLayer({ id: 'facilities-circle', type: 'circle', source: 'market-facilities', filter: ['==', ['geometry-type'], 'Point'], paint: { 'circle-radius': 8, 'circle-color': '#ffffff', 'circle-stroke-color': '#2563eb', 'circle-stroke-width': 3 } })
  map.addLayer({ id: 'facilities-label', type: 'symbol', source: 'market-facilities', filter: ['==', ['geometry-type'], 'Point'], layout: { 'text-field': ['get', 'name'], 'text-size': 12, 'text-offset': [0, 1.5], 'text-anchor': 'top' }, paint: { 'text-color': '#1e3a8a', 'text-halo-color': '#ffffff', 'text-halo-width': 2 } })
  map.addLayer({ id: 'zone-drawing-fill', type: 'fill', source: 'zone-drawing', filter: ['==', ['geometry-type'], 'Polygon'], paint: { 'fill-color': '#2563eb', 'fill-opacity': .16 } })
  map.addLayer({ id: 'zone-drawing-line', type: 'line', source: 'zone-drawing', paint: { 'line-color': '#2563eb', 'line-width': 3, 'line-dasharray': [2, 1] } })
  map.addLayer({ id: 'zone-drawing-points', type: 'circle', source: 'zone-drawing', filter: ['==', ['geometry-type'], 'Point'], paint: { 'circle-radius': 5, 'circle-color': '#2563eb', 'circle-stroke-color': '#ffffff', 'circle-stroke-width': 2 } })

  // ขอบขาวรองใต้เส้นกรอบ ให้คมชัดบนแผนที่ถนน
  map.addLayer({
    id: 'projects-casing',
    type: 'line',
    source: 'projects-polygons',
    layout: { 'line-join': 'round' },
    paint: {
      'line-color': '#ffffff',
      'line-width': ['case', ['boolean', ['get', 'isActive'], false], 8, ['boolean', ['get', 'isHovered'], false], 8, 5],
      'line-opacity': ['case', ['boolean', ['get', 'isHovered'], false], 1, ['boolean', ['get', 'isDimmed'], false], 0.3, 0.95]
    }
  })

  // เส้นกรอบโครงการ
  map.addLayer({
    id: 'projects-stroke',
    type: 'line',
    source: 'projects-polygons',
    layout: { 'line-join': 'round' },
    paint: {
      'line-color': ['get', 'strokeColor'],
      'line-width': ['case', ['boolean', ['get', 'isActive'], false], 4.5, ['boolean', ['get', 'isHovered'], false], 4.5, 2.5],
      'line-opacity': ['case', ['boolean', ['get', 'isHovered'], false], 1, ['boolean', ['get', 'isDimmed'], false], 0.3, 1]
    }
  })

  // แผงย่อยภายในโครงการ (เน้นแผงที่เลือกให้เข้มเด่น แผงอื่นจางลง)
  map.addLayer({
    id: 'stalls-fill',
    type: 'fill',
    source: 'project-stalls',
    paint: {
      'fill-color': [
        'case',
        ['boolean', ['get', 'isSelected'], false], '#2563eb',
        ['boolean', ['get', 'isDimmed'], false], [
          'case',
          ['==', ['get', 'status'], 'available'], '#86efac',
          ['==', ['get', 'status'], 'reserved'], '#fed7aa',
          '#94a3b8'
        ],
        ['==', ['get', 'status'], 'available'], '#10b981',
        ['==', ['get', 'status'], 'reserved'], '#f59e0b',
        '#475569'
      ],
      'fill-opacity': [
        'case',
        ['boolean', ['get', 'isSelected'], false], 1,
        ['boolean', ['get', 'isDimmed'], false], 0.32,
        1
      ],
      'fill-outline-color': '#ffffff'
    }
  })

  // เส้นแบ่งล็อคแผงสถาปัตยกรรมสีขาวสะอาดตา คมชัดระดับพรีเมียม (Clean White Architectural Partitions)
  map.addLayer({
    id: 'stalls-line',
    type: 'line',
    source: 'project-stalls',
    layout: {
      'line-join': 'round',
      'line-cap': 'round'
    },
    paint: {
      'line-color': '#ffffff',
      'line-width': [
        'interpolate', ['linear'], ['zoom'],
        14, 0.8,
        16, 1.4,
        18, 2.0,
        20, 2.8
      ],
      'line-opacity': [
        'case',
        ['boolean', ['get', 'isDimmed'], false], 0.45,
        1
      ]
    }
  })

  // Show the survey stall number inside every polygon once the cells are
  // large enough to read. Internal booking codes remain unchanged.
  map.addLayer({
    id: 'stalls-label',
    type: 'symbol',
    source: 'project-stalls',
    minzoom: 18,
    filter: ['!=', ['get', 'label'], ''],
    layout: {
      'text-field': ['get', 'label'],
      'visibility': showStallNumberLabels.value ? 'visible' : 'none',
      'text-size': [
        'interpolate', ['linear'], ['zoom'],
        18, 8,
        20, 11,
        22, 13
      ],
      'text-anchor': 'center',
      'text-allow-overlap': true,
      'text-ignore-placement': true,
      'text-rotation-alignment': 'viewport',
      'text-pitch-alignment': 'viewport'
    },
    paint: {
      'text-color': [
        'case',
        ['boolean', ['get', 'isSelected'], false], '#ffffff',
        '#0f172a'
      ],
      'text-halo-color': [
        'case',
        ['boolean', ['get', 'isSelected'], false], '#1d4ed8',
        '#ffffff'
      ],
      'text-halo-width': 1.25,
      'text-opacity': [
        'case',
        ['boolean', ['get', 'isDimmed'], false], 0.55,
        0.95
      ]
    }
  })

  // Search matches keep their real status colour and receive a clear blue
  // outline so users can distinguish the target from dimmed surrounding stalls.
  map.addLayer({
    id: 'stalls-search-casing',
    type: 'line',
    source: 'project-stalls',
    filter: ['boolean', ['get', 'isSearchMatch'], false],
    layout: { 'line-join': 'round', 'line-cap': 'round' },
    paint: {
      'line-color': '#ffffff',
      'line-width': ['interpolate', ['linear'], ['zoom'], 14, 3.6, 16, 4.8, 18, 6.2, 20, 8],
      'line-opacity': 1
    }
  })

  map.addLayer({
    id: 'stalls-search-stroke',
    type: 'line',
    source: 'project-stalls',
    filter: ['boolean', ['get', 'isSearchMatch'], false],
    layout: { 'line-join': 'round', 'line-cap': 'round' },
    paint: {
      'line-color': '#2563eb',
      'line-width': ['interpolate', ['linear'], ['zoom'], 14, 1.8, 16, 2.4, 18, 3.2, 20, 4.2],
      'line-opacity': 1
    }
  })

  // Keep the real status colour visible and identify the exact stall being
  // inspected with a white halo and a blue focus outline.
  map.addLayer({
    id: 'stalls-hover-casing',
    type: 'line',
    source: 'project-stalls',
    layout: {
      'line-join': 'round',
      'line-cap': 'round'
    },
    paint: {
      'line-color': '#ffffff',
      'line-width': [
        'interpolate', ['linear'], ['zoom'],
        14, 2.4,
        16, 3.2,
        18, 4,
        20, 5
      ],
      'line-opacity': ['case', ['boolean', ['feature-state', 'isHovered'], false], 1, 0]
    }
  })

  map.addLayer({
    id: 'stalls-hover-stroke',
    type: 'line',
    source: 'project-stalls',
    layout: {
      'line-join': 'round',
      'line-cap': 'round'
    },
    paint: {
      'line-color': '#2563eb',
      'line-width': [
        'interpolate', ['linear'], ['zoom'],
        14, 1.2,
        16, 1.6,
        18, 2.2,
        20, 3
      ],
      'line-opacity': ['case', ['boolean', ['feature-state', 'isHovered'], false], 1, 0]
    }
  })

  // ขอบสีขาวรองใต้กรอบแผงที่เลือก เพื่อให้ยังเห็นตำแหน่งชัดบนทุกสีของผัง
  map.addLayer({
    id: 'stalls-selected-casing',
    type: 'line',
    source: 'project-stalls',
    filter: ['boolean', ['get', 'isSelected'], false],
    layout: {
      'line-join': 'round',
      'line-cap': 'round'
    },
    paint: {
      'line-color': '#ffffff',
      'line-width': [
        'interpolate', ['linear'], ['zoom'],
        14, 4.8,
        16, 6.2,
        18, 8,
        20, 10
      ],
      'line-opacity': 1
    }
  })

  // เส้นขอบแผงที่ถูกเลือก (Selected Stall Stroke) - สีน้ำเงินสดเด่นชัด
  map.addLayer({
    id: 'stalls-selected-stroke',
    type: 'line',
    source: 'project-stalls',
    filter: ['boolean', ['get', 'isSelected'], false],
    layout: {
      'line-join': 'round',
      'line-cap': 'round'
    },
    paint: {
      'line-color': '#1d4ed8',
      'line-width': [
        'interpolate', ['linear'], ['zoom'],
        14, 2.2,
        16, 3.2,
        18, 4.4,
        20, 5.8
      ],
      'line-opacity': 1
    }
  })

  // Imported market stalls often touch the project boundary. Their white
  // separators would otherwise paint over the blue outline and make it look
  // broken. Redraw only the active project's outline above every stall layer.
  map.addLayer({
    id: 'projects-active-outline',
    type: 'line',
    source: 'projects-polygons',
    filter: ['==', ['get', 'isActive'], true],
    layout: {
      'line-join': 'round',
      'line-cap': 'round'
    },
    paint: {
      'line-color': '#2563eb',
      'line-width': [
        'interpolate', ['linear'], ['zoom'],
        14, 2.4,
        16, 3.2,
        18, 4,
        20, 5
      ],
      'line-opacity': 1
    }
  })

  // Project polygons live below the active project's stall layers. Draw the
  // hovered project once more on top so overlapping zones still visibly react
  // together with their number badge.
  map.addLayer({
    id: 'projects-hover-casing',
    type: 'line',
    source: 'projects-polygons',
    filter: ['==', ['get', 'isHovered'], true],
    layout: {
      'line-join': 'round',
      'line-cap': 'round'
    },
    paint: {
      'line-color': '#ffffff',
      'line-width': [
        'interpolate', ['linear'], ['zoom'],
        14, 6,
        16, 8,
        18, 10,
        20, 12
      ],
      'line-opacity': 1
    }
  })

  map.addLayer({
    id: 'projects-hover-outline',
    type: 'line',
    source: 'projects-polygons',
    filter: ['==', ['get', 'isHovered'], true],
    layout: {
      'line-join': 'round',
      'line-cap': 'round'
    },
    paint: {
      'line-color': '#2563eb',
      'line-width': [
        'interpolate', ['linear'], ['zoom'],
        14, 3,
        16, 4,
        18, 5.5,
        20, 7
      ],
      'line-opacity': 1
    }
  })

  // อาคาร 3D (แทรกใต้ projects-fill)
  add3DBuildings()

  renderZoneLabels()
}

const refreshGeometrySources = () => {
  if (!map) return
  map.getSource('projects-polygons')?.setData(getProjectsGeoJson())
  map.getSource('project-stalls')?.setData(getStallsGeoJson())
  map.getSource('market-facilities')?.setData(getFacilitiesGeoJson())
  refreshDrawingZone()
}

const refreshSources = () => {
  if (!map) return
  refreshGeometrySources()
  renderZoneLabels()
  renderFacilityMarkers()
}

// Availability updates do not change facilities or drawing geometry. Updating
// those sources and rebuilding their DOM markers made contract data arriving
// during a fly/fit animation visibly stutter on the public booking map.
let statusRefreshFrame: number | null = null
let statusRefreshQueuedForMove = false
const refreshStatusSources = () => {
  if (!map) return
  map.getSource('projects-polygons')?.setData(getProjectsGeoJson())
  map.getSource('project-stalls')?.setData(getStallsGeoJson())
  renderZoneLabels()
}
const scheduleStatusSourceRefresh = () => {
  if (!map || isUnmounted) return
  if (map.isMoving()) {
    if (statusRefreshQueuedForMove) return
    statusRefreshQueuedForMove = true
    map.once('moveend', () => {
      statusRefreshQueuedForMove = false
      scheduleStatusSourceRefresh()
    })
    return
  }
  if (statusRefreshFrame !== null) window.cancelAnimationFrame(statusRefreshFrame)
  statusRefreshFrame = window.requestAnimationFrame(() => {
    statusRefreshFrame = null
    if (!isUnmounted) refreshStatusSources()
  })
}

const clearPolygonVertexHandles = () => {
  polygonHandleMarkers.forEach(marker => {
    try { marker.remove() } catch {}
  })
  polygonHandleMarkers = []
}

const clearStallVertexHandles = () => {
  stallHandleMarkers.forEach(marker => {
    try { marker.remove() } catch {}
  })
  stallHandleMarkers = []
}

const calcPolygonMetrics = (ring: [number, number][]) => {
  if (ring.length < 3) return { areaSqm: 4, widthM: 2, depthM: 2, centroid: ring[0] || ([100.5, 13.8] as [number, number]) }
  const open = ring.length > 1 && ring[0]?.[0] === ring.at(-1)?.[0] && ring[0]?.[1] === ring.at(-1)?.[1]
    ? ring.slice(0, -1)
    : ring
  const latRad = ((open[0]?.[1] || 13.8) * Math.PI) / 180
  const mPerDegLng = 111320 * Math.cos(latRad)
  const mPerDegLat = 110540
  let twiceArea = 0
  for (let i = 0; i < ring.length - 1; i++) {
    const x1 = (ring[i]?.[0] || 0) * mPerDegLng
    const y1 = (ring[i]?.[1] || 0) * mPerDegLat
    const x2 = (ring[i + 1]?.[0] || 0) * mPerDegLng
    const y2 = (ring[i + 1]?.[1] || 0) * mPerDegLat
    twiceArea += x1 * y2 - x2 * y1
  }
  const areaSqm = Math.max(1, Number((Math.abs(twiceArea) / 2).toFixed(1)))
  const centroid: [number, number] = [
    open.reduce((s, p) => s + p[0], 0) / open.length,
    open.reduce((s, p) => s + p[1], 0) / open.length
  ]
  if (open.length >= 3) {
    const d1 = Math.hypot((open[1]![0] - open[0]![0]) * mPerDegLng, (open[1]![1] - open[0]![1]) * mPerDegLat)
    const d2 = Math.hypot((open[2]![0] - open[1]![0]) * mPerDegLng, (open[2]![1] - open[1]![1]) * mPerDegLat)
    return {
      areaSqm,
      widthM: Math.max(0.5, Number(d1.toFixed(1))),
      depthM: Math.max(0.5, Number(d2.toFixed(1))),
      centroid
    }
  }
  return { areaSqm, widthM: 2, depthM: 2, centroid }
}

// Staff can adjust a real project boundary directly on the map (both from Overview and inside a Zone).
const renderPolygonVertexHandles = () => {
  clearPolygonVertexHandles()
  const targetProject = polygonEditTargetProject.value
  if (!map || !maplibregl || !props.managementMode || !isPolygonEditMode.value || !targetProject) return

  const projectID = targetProject.id
  const sourceProject = currentMarket.value.projects.find(item => item.id === projectID)
  if (!sourceProject) return
  const closedCoordinates = sourceProject.coordinates || []
  const vertexCount = closedCoordinates.length > 1 && closedCoordinates[0]?.[0] === closedCoordinates.at(-1)?.[0] && closedCoordinates[0]?.[1] === closedCoordinates.at(-1)?.[1]
    ? closedCoordinates.length - 1
    : closedCoordinates.length

  const midpointMarkersByEdge: any[] = []

  // 1. Midpoint '+' handles along each edge to insert & drag a new corner vertex
  for (let index = 0; index < vertexCount; index++) {
    const p1 = closedCoordinates[index]
    const p2 = closedCoordinates[(index + 1) % vertexCount]
    if (!p1 || !p2) continue
    const midPt: [number, number] = [(p1[0] + p2[0]) / 2, (p1[1] + p2[1]) / 2]

    const midEl = document.createElement('button')
    midEl.type = 'button'
    midEl.className = 'zone-midpoint-handle'
    midEl.setAttribute('aria-label', `เพิ่มจุดมุมโซนระหว่างจุดที่ ${index + 1} และ ${((index + 1) % vertexCount) + 1}`)
    midEl.title = 'คลิกหรือลากเพื่อเพิ่มจุดมุมโซนใหม่บนเส้นขอบนี้'
    Object.assign(midEl.style, {
      width: '16px', height: '16px', borderRadius: '50%', border: '2px solid #2563eb',
      background: '#ffffff', color: '#2563eb', boxShadow: '0 1px 4px rgba(15,23,42,.3)',
      cursor: 'copy', padding: '0', display: 'flex', alignItems: 'center', justifyContent: 'center',
      zIndex: '24'
    })
    midEl.innerHTML = '<svg width="10" height="10" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3" stroke-linecap="round"><path d="M12 5v14M5 12h14"/></svg>'

    const midMarker = new maplibregl.Marker({ element: midEl, draggable: true, anchor: 'center' })
      .setLngLat(midPt)
      .addTo(map)

    let insertedIdx = -1
    let didDrag = false

    const insertVertexAt = (lngLat: [number, number]): number => {
      const target = currentMarket.value.projects.find(item => item.id === projectID)
      if (!target) return -1
      const openVerts = target.coordinates.slice(0, vertexCount).map(pair => [...pair] as [number, number])
      const insertPos = index + 1
      openVerts.splice(insertPos, 0, [lngLat[0], lngLat[1]])
      const nextClosed: [number, number][] = [...openVerts, [...openVerts[0]!]]
      const metrics = calcPolygonMetrics(nextClosed)
      target.coordinates = nextClosed
      target.realAreaSqm = Math.round(metrics.areaSqm)
      target.centroid = metrics.centroid
      clearLayoutCache()
      refreshSources()
      return insertPos
    }

    midMarker.on('dragstart', () => {
      didDrag = true
      rememberPolygonEdit()
      const cur = midMarker.getLngLat()
      insertedIdx = insertVertexAt([cur.lng, cur.lat])
      midEl.style.background = '#3b5bff'
      midEl.style.borderColor = '#ffffff'
      midEl.style.color = '#ffffff'
      midEl.style.cursor = 'grabbing'
    })

    midMarker.on('drag', () => {
      if (insertedIdx < 0) return
      const next = midMarker.getLngLat()
      const target = currentMarket.value.projects.find(item => item.id === projectID)
      if (!target) return
      const coordinates = target.coordinates.map(pair => [...pair] as [number, number])
      coordinates[insertedIdx] = [next.lng, next.lat]
      if (insertedIdx === 0 && coordinates.length > 1) {
        coordinates[coordinates.length - 1] = [next.lng, next.lat]
      }
      const metrics = calcPolygonMetrics(coordinates)
      target.coordinates = coordinates
      target.realAreaSqm = Math.round(metrics.areaSqm)
      target.centroid = metrics.centroid
      clearLayoutCache()
      refreshSources()
    })

    midMarker.on('dragend', () => {
      const target = currentMarket.value.projects.find(item => item.id === projectID)
      if (!target) return
      const coordinates = target.coordinates.map(pair => [...pair] as [number, number])
      const metrics = calcPolygonMetrics(coordinates)
      emit('updateProjectPolygon', {
        ...target,
        coordinates,
        realAreaSqm: Math.round(metrics.areaSqm),
        centroid: metrics.centroid
      })
      renderPolygonVertexHandles()
    })

    midEl.addEventListener('click', (ev) => {
      ev.stopPropagation()
      if (didDrag) {
        didDrag = false
        return
      }
      const cur = midMarker.getLngLat()
      rememberPolygonEdit()
      if (insertVertexAt([cur.lng, cur.lat]) >= 0) {
        const target = currentMarket.value.projects.find(item => item.id === projectID)
        if (target) {
          emit('updateProjectPolygon', {
            ...target,
            coordinates: target.coordinates.map(pair => [...pair] as [number, number]),
            realAreaSqm: target.realAreaSqm,
            centroid: target.centroid
          })
        }
        renderPolygonVertexHandles()
      }
    })

    midpointMarkersByEdge.push(midMarker)
    polygonHandleMarkers.push(midMarker)
  }

  // 2. Corner vertex handles (drag to move, double-click / right-click to delete when > 3 vertices)
  for (let index = 0; index < vertexCount; index++) {
    const point = closedCoordinates[index]
    if (!point) continue
    const element = document.createElement('button')
    element.type = 'button'
    element.className = 'zone-vertex-handle'
    element.setAttribute('aria-label', `ปรับจุดขอบเขต ${index + 1}`)
    element.title = vertexCount > 3
      ? `ลากเพื่อปรับจุดขอบเขตโซนที่ ${index + 1} (ดับเบิลคลิกเพื่อลบมุมนี้)`
      : `ลากเพื่อปรับจุดขอบเขตโซนที่ ${index + 1}`
    Object.assign(element.style, {
      width: '18px', height: '18px', borderRadius: '50%', border: '3px solid #ffffff',
      background: '#3b5bff', boxShadow: '0 1px 5px rgba(15,23,42,.35)', cursor: 'grab', padding: '0',
      zIndex: '25'
    })
    const marker = new maplibregl.Marker({ element, draggable: true, anchor: 'center' })
      .setLngLat(point)
      .addTo(map)
    marker.on('dragstart', () => {
      rememberPolygonEdit()
      element.style.cursor = 'grabbing'
    })
    marker.on('drag', () => {
      const next = marker.getLngLat()
      const target = currentMarket.value.projects.find(item => item.id === projectID)
      if (!target) return
      const coordinates = target.coordinates.map(pair => [...pair] as [number, number])
      coordinates[index] = [next.lng, next.lat]
      if (coordinates.length > vertexCount) coordinates[coordinates.length - 1] = [...coordinates[0]!]
      const metrics = calcPolygonMetrics(coordinates)
      target.coordinates = coordinates
      target.realAreaSqm = Math.round(metrics.areaSqm)
      target.centroid = metrics.centroid
      // Keep adjacent midpoint handles centered while dragging
      const prevIdx = (index - 1 + vertexCount) % vertexCount
      const nextIdx = (index + 1) % vertexCount
      const prevPt = coordinates[prevIdx]
      const nextPt = coordinates[nextIdx]
      if (prevPt && midpointMarkersByEdge[prevIdx]) {
        midpointMarkersByEdge[prevIdx].setLngLat([(prevPt[0] + next.lng) / 2, (prevPt[1] + next.lat) / 2])
      }
      if (nextPt && midpointMarkersByEdge[index]) {
        midpointMarkersByEdge[index].setLngLat([(next.lng + nextPt[0]) / 2, (next.lat + nextPt[1]) / 2])
      }
      clearLayoutCache()
      refreshSources()
    })
    marker.on('dragend', () => {
      element.style.cursor = 'grab'
      const target = currentMarket.value.projects.find(item => item.id === projectID)
      if (!target) return
      const coordinates = target.coordinates.map(pair => [...pair] as [number, number])
      const metrics = calcPolygonMetrics(coordinates)
      emit('updateProjectPolygon', {
        ...target,
        coordinates,
        realAreaSqm: Math.round(metrics.areaSqm),
        centroid: metrics.centroid
      })
      renderPolygonVertexHandles()
    })
    const removeCornerVertex = (ev: Event) => {
      ev.preventDefault()
      ev.stopPropagation()
      if (vertexCount <= 3) return
      rememberPolygonEdit()
      const target = currentMarket.value.projects.find(item => item.id === projectID)
      if (!target) return
      const openVerts = target.coordinates.slice(0, vertexCount).map(pair => [...pair] as [number, number])
      openVerts.splice(index, 1)
      if (openVerts.length < 3) return
      const nextClosed: [number, number][] = [...openVerts, [...openVerts[0]!]]
      const metrics = calcPolygonMetrics(nextClosed)
      target.coordinates = nextClosed
      target.realAreaSqm = Math.round(metrics.areaSqm)
      target.centroid = metrics.centroid
      clearLayoutCache()
      refreshSources()
      emit('updateProjectPolygon', {
        ...target,
        coordinates: nextClosed,
        realAreaSqm: Math.round(metrics.areaSqm),
        centroid: metrics.centroid
      })
      renderPolygonVertexHandles()
    }
    element.addEventListener('dblclick', removeCornerVertex)
    element.addEventListener('contextmenu', removeCornerVertex)
    polygonHandleMarkers.push(marker)
  }
}

const renderStallVertexHandles = () => {
  clearStallVertexHandles()
  if (!map || !maplibregl || !props.managementMode || !isStallPolygonEditMode.value || !activeProject.value) return

  const projectID = activeProject.value.id
  const sourceProject = currentMarket.value.projects.find(item => item.id === projectID)
  if (!sourceProject) return

  // Ensure stalls array is materialized if editing a procedurally generated zone
  if (!sourceProject.stalls || sourceProject.stalls.length === 0) {
    const allGenerated: any[] = []
    const floors = sourceProject.floors?.length ? sourceProject.floors : [{ id: '1', name: 'ชั้น 1', shortLabel: '1F', level: 1 }]
    floors.forEach(fl => {
      const flLayout = applyMarketStatus(getProjectStallLayout(sourceProject, fl.id))
      flLayout.stalls.forEach(st => {
        allGenerated.push({ ...st, floor: String(fl.id) })
      })
    })
    sourceProject.stalls = cloneProjects(allGenerated)
    sourceProject.layoutMode = 'manual'
    clearLayoutCache()
  }

  const targetCode = editingStallCode.value || activeLayout.value?.stalls[0]?.code
  if (!targetCode) return
  editingStallCode.value = targetCode

  const stallItem = (sourceProject.stalls || []).find(
    (s: any) => String(s.code) === String(targetCode) && String(s.floor || '1') === String(currentFloorId.value)
  ) || (sourceProject.stalls || []).find((s: any) => String(s.code) === String(targetCode))
  if (!stallItem || !Array.isArray(stallItem.coordinates)) return

  const closedCoordinates: [number, number][] = stallItem.coordinates
  const vertexCount = closedCoordinates.length > 1 && closedCoordinates[0]?.[0] === closedCoordinates.at(-1)?.[0] && closedCoordinates[0]?.[1] === closedCoordinates.at(-1)?.[1]
    ? closedCoordinates.length - 1
    : closedCoordinates.length

  // 1. Midpoint '+' handles along each stall edge
  for (let index = 0; index < vertexCount; index++) {
    const p1 = closedCoordinates[index]
    const p2 = closedCoordinates[(index + 1) % vertexCount]
    if (!p1 || !p2) continue
    const midPt: [number, number] = [(p1[0] + p2[0]) / 2, (p1[1] + p2[1]) / 2]

    const midEl = document.createElement('button')
    midEl.type = 'button'
    midEl.className = 'stall-midpoint-handle'
    midEl.setAttribute('aria-label', `เพิ่มจุดมุมแผง ${targetCode} บนเส้นขอบที่ ${index + 1}`)
    midEl.title = `คลิกหรือลากเพื่อเพิ่มจุดมุมแผง ${targetCode}`
    Object.assign(midEl.style, {
      width: '14px', height: '14px', borderRadius: '50%', border: '2px solid #059669',
      background: '#ffffff', color: '#059669', boxShadow: '0 1px 4px rgba(15,23,42,.3)',
      cursor: 'copy', padding: '0', display: 'flex', alignItems: 'center', justifyContent: 'center',
      zIndex: '24'
    })
    midEl.innerHTML = '<svg width="8" height="8" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3" stroke-linecap="round"><path d="M12 5v14M5 12h14"/></svg>'

    const midMarker = new maplibregl.Marker({ element: midEl, draggable: true, anchor: 'center' })
      .setLngLat(midPt)
      .addTo(map)

    let insertedIdx = -1
    let didDrag = false

    const insertStallVertexAt = (lngLat: [number, number]): number => {
      const target = currentMarket.value.projects.find(item => item.id === projectID)
      if (!target || !Array.isArray(target.stalls)) return -1
      const sTarget = target.stalls.find((s: any) => String(s.code) === String(targetCode))
      if (!sTarget) return -1
      const openVerts = sTarget.coordinates.slice(0, vertexCount).map((pair: [number, number]) => [...pair] as [number, number])
      const insertPos = index + 1
      openVerts.splice(insertPos, 0, [lngLat[0], lngLat[1]])
      const nextClosed: [number, number][] = [...openVerts, [...openVerts[0]!]]
      const metrics = calcPolygonMetrics(nextClosed)
      sTarget.coordinates = nextClosed
      sTarget.widthM = metrics.widthM
      sTarget.depthM = metrics.depthM
      sTarget.size = metrics.areaSqm
      sTarget.center = metrics.centroid
      clearLayoutCache()
      refreshSources()
      return insertPos
    }

    midMarker.on('dragstart', () => {
      didDrag = true
      const cur = midMarker.getLngLat()
      insertedIdx = insertStallVertexAt([cur.lng, cur.lat])
      midEl.style.background = '#10b981'
      midEl.style.borderColor = '#ffffff'
      midEl.style.color = '#ffffff'
      midEl.style.cursor = 'grabbing'
    })

    midMarker.on('drag', () => {
      if (insertedIdx < 0) return
      const next = midMarker.getLngLat()
      const target = currentMarket.value.projects.find(item => item.id === projectID)
      if (!target || !Array.isArray(target.stalls)) return
      const sTarget = target.stalls.find((s: any) => String(s.code) === String(targetCode))
      if (!sTarget) return
      const coords = sTarget.coordinates.map((pair: [number, number]) => [...pair] as [number, number])
      coords[insertedIdx] = [next.lng, next.lat]
      if (insertedIdx === 0 && coords.length > 1) {
        coords[coords.length - 1] = [next.lng, next.lat]
      }
      const metrics = calcPolygonMetrics(coords)
      sTarget.coordinates = coords
      sTarget.widthM = metrics.widthM
      sTarget.depthM = metrics.depthM
      sTarget.size = metrics.areaSqm
      sTarget.center = metrics.centroid
      clearLayoutCache()
      refreshSources()
    })

    midMarker.on('dragend', () => {
      emitCurrentEditingStallUpdate()
      renderStallVertexHandles()
    })

    midEl.addEventListener('click', (ev) => {
      ev.stopPropagation()
      if (didDrag) {
        didDrag = false
        return
      }
      const cur = midMarker.getLngLat()
      rememberStallEdit()
      if (insertStallVertexAt([cur.lng, cur.lat]) >= 0) {
        emitCurrentEditingStallUpdate()
        renderStallVertexHandles()
      }
    })

    stallHandleMarkers.push(midMarker)
  }

  // 2. Corner vertex handles
  for (let index = 0; index < vertexCount; index++) {
    const point = closedCoordinates[index]
    if (!point) continue
    const element = document.createElement('button')
    element.type = 'button'
    element.className = 'stall-vertex-handle'
    element.setAttribute('aria-label', `ปรับจุดมุมแผง ${targetCode} จุดที่ ${index + 1}`)
    element.title = `ลากเพื่อปรับขนาดแผง ${targetCode} (มุมที่ ${index + 1})`
    Object.assign(element.style, {
      width: '16px', height: '16px', borderRadius: '50%', border: '2.5px solid #ffffff',
      background: '#10b981', boxShadow: '0 1px 5px rgba(15,23,42,.38)', cursor: 'grab', padding: '0',
      zIndex: '25'
    })
    const marker = new maplibregl.Marker({ element, draggable: true, anchor: 'center' })
      .setLngLat(point)
      .addTo(map)
    marker.on('dragstart', () => {
      rememberStallEdit()
      element.style.cursor = 'grabbing'
    })
    marker.on('drag', () => {
      const next = marker.getLngLat()
      const target = currentMarket.value.projects.find(item => item.id === projectID)
      if (!target || !Array.isArray(target.stalls)) return
      const sTarget = target.stalls.find((s: any) => String(s.code) === String(targetCode))
      if (!sTarget) return
      const coords = sTarget.coordinates.map((pair: [number, number]) => [...pair] as [number, number])
      coords[index] = [next.lng, next.lat]
      if (coords.length > vertexCount) coords[coords.length - 1] = [...coords[0]!]
      const metrics = calcPolygonMetrics(coords)
      sTarget.coordinates = coords
      sTarget.widthM = metrics.widthM
      sTarget.depthM = metrics.depthM
      sTarget.size = metrics.areaSqm
      sTarget.center = metrics.centroid
      clearLayoutCache()
      refreshSources()
    })
    marker.on('dragend', () => {
      element.style.cursor = 'grab'
      emitCurrentEditingStallUpdate()
      renderStallVertexHandles()
    })
    stallHandleMarkers.push(marker)
  }
}

watch(() => props.inspectOnly, (inspect) => {
  if (!inspect) return
  cancelDrawing()
  isPolygonEditMode.value = false
  isStallPolygonEditMode.value = false
  editingStallCode.value = null
  clearPolygonVertexHandles()
  clearStallVertexHandles()
})

const FACILITY_META: Record<string, { label: string; color: string; bg: string; svg: string }> = {
  sanitation: {
    label: 'สุขาภิบาล',
    color: '#0284c7',
    bg: '#e0f2fe',
    svg: '<svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.2" stroke-linecap="round" stroke-linejoin="round"><path d="M12 2v20M5 7h14M7 7v13M17 7v13"/><path d="M9 11h6M9 15h6"/></svg>'
  },
  access: {
    label: 'ทางสัญจร',
    color: '#d97706',
    bg: '#fef3c7',
    svg: '<svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.2" stroke-linecap="round" stroke-linejoin="round"><path d="M5 12h14M13 6l6 6-6 6"/></svg>'
  },
  'common-area': {
    label: 'พื้นที่ส่วนกลาง',
    color: '#0f766e',
    bg: '#ccfbf1',
    svg: '<svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.2" stroke-linecap="round" stroke-linejoin="round"><path d="M3 21h18M5 21V8l7-5 7 5v13"/><path d="M9 21v-6h6v6"/></svg>'
  },
  building: {
    label: 'อาคารและสิ่งปลูกสร้าง',
    color: '#64748b',
    bg: '#f1f5f9',
    svg: '<svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.2" stroke-linecap="round" stroke-linejoin="round"><rect x="4" y="3" width="16" height="18" rx="1"/><path d="M8 7h2M14 7h2M8 11h2M14 11h2M9 21v-5h6v5"/></svg>'
  },
  toilet: {
    label: 'ห้องน้ำ',
    color: '#0284c7',
    bg: '#e0f2fe',
    svg: '<svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.1" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true"><circle cx="6" cy="4" r="1.8"/><path d="M3 21v-7H2v-4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v4H9v7M6 14v7"/><path d="M12 3v18"/><circle cx="18" cy="4" r="1.8"/><path d="M16 21v-4h-2l2-9h4l2 9h-2v4M18 17v4"/></svg>'
  },
  parking: {
    label: 'ที่จอดรถ',
    color: '#2563eb',
    bg: '#dbeafe',
    svg: '<svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.2" stroke-linecap="round" stroke-linejoin="round"><rect width="18" height="18" x="3" y="3" rx="3"/><path d="M9 17V7h4a3 3 0 0 1 0 6H9"/></svg>'
  },
  exit: {
    label: 'ทางเข้า–ออก',
    color: '#059669',
    bg: '#d1fae5',
    svg: '<svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.2" stroke-linecap="round" stroke-linejoin="round"><path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h8"/><path d="M10 12h11m0 0-4-4m4 4-4 4"/></svg>'
  },
  service: {
    label: 'จุดบริการประชาสัมพันธ์',
    color: '#7c3aed',
    bg: '#ede9fe',
    svg: '<svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.2" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="10"/><path d="M12 16v-4M12 8h.01"/></svg>'
  },
  waste: {
    label: 'จุดทิ้งขยะ',
    color: '#d97706',
    bg: '#fef3c7',
    svg: '<svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.2" stroke-linecap="round" stroke-linejoin="round"><path d="M3 6h18M8 6V4h8v2M19 6l-1 14H6L5 6M10 11v5M14 11v5"/></svg>'
  }
}

const getFacilityMeta = (type: string) => FACILITY_META[type] || {
  label: 'จุดอำนวยความสะดวก',
  color: '#2563eb',
  bg: '#eff6ff',
  svg: '<svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.2" stroke-linecap="round" stroke-linejoin="round"><path d="M12 22s7-5.3 7-12a7 7 0 1 0-14 0c0 6.7 7 12 7 12Z"/><circle cx="12" cy="10" r="2.5"/></svg>'
}

const NAMED_FACILITY_META = {
  sink: {
    label: 'อ่างล้างมือ', color: '#0284c7', bg: '#e0f2fe',
    svg: '<svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.2" stroke-linecap="round" stroke-linejoin="round"><path d="M5 12h14v2a5 5 0 0 1-5 5h-4a5 5 0 0 1-5-5v-2Z"/><path d="M8 12V8a4 4 0 0 1 8 0v1M14 9h4M7 19v2M17 19v2"/></svg>'
  },
  stairs: {
    label: 'บันได', color: '#d97706', bg: '#fef3c7',
    svg: '<svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.2" stroke-linecap="round" stroke-linejoin="round"><path d="M3 19h5v-4h4v-4h4V7h5"/><path d="M3 22h18"/></svg>'
  },
  escalator: {
    label: 'บันไดเลื่อน', color: '#d97706', bg: '#fef3c7',
    svg: '<svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.2" stroke-linecap="round" stroke-linejoin="round"><path d="M4 18h4l8-12h4M4 21h5l8-12h3"/><path d="m16 17 4 4m0-4v4h-4"/></svg>'
  },
  shrine: {
    label: 'ศาล', color: '#dc2626', bg: '#fee2e2',
    svg: '<svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.2" stroke-linecap="round" stroke-linejoin="round"><path d="M3 9h18L12 3 3 9Z"/><path d="M5 9v9M19 9v9M9 9v9M15 9v9M3 21h18M7 18h10"/></svg>'
  },
  clockTower: {
    label: 'หอนาฬิกา', color: '#7c3aed', bg: '#ede9fe',
    svg: '<svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.2" stroke-linecap="round" stroke-linejoin="round"><path d="M7 22h10L16 9H8L7 22Z"/><circle cx="12" cy="6" r="4"/><path d="M12 4v2l1.5 1M5 22h14"/></svg>'
  },
  tent: {
    label: 'เต็นท์', color: '#059669', bg: '#d1fae5',
    svg: '<svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.2" stroke-linecap="round" stroke-linejoin="round"><path d="M3 20 12 4l9 16H3Z"/><path d="m9 20 3-6 3 6M12 4v10"/></svg>'
  },
  plaza: {
    label: 'ลานส่วนกลาง', color: '#0f766e', bg: '#ccfbf1',
    svg: '<svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.2" stroke-linecap="round" stroke-linejoin="round"><rect x="3" y="3" width="18" height="18" rx="2"/><path d="M8 8h8v8H8zM3 12h5M16 12h5M12 3v5M12 16v5"/></svg>'
  },
  meeting: {
    label: 'ห้องประชุม', color: '#7c3aed', bg: '#ede9fe',
    svg: '<svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.2" stroke-linecap="round" stroke-linejoin="round"><circle cx="8" cy="8" r="2"/><circle cx="16" cy="8" r="2"/><path d="M4 19v-2a4 4 0 0 1 8 0v2M12 17a4 4 0 0 1 8 0v2"/></svg>'
  },
  training: {
    label: 'ศูนย์ฝึกอาชีพ', color: '#2563eb', bg: '#dbeafe',
    svg: '<svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.2" stroke-linecap="round" stroke-linejoin="round"><path d="m3 9 9-5 9 5-9 5-9-5Z"/><path d="M7 12v5c3 2 7 2 10 0v-5M21 9v6"/></svg>'
  },
  office: {
    label: 'สำนักงาน', color: '#64748b', bg: '#f1f5f9',
    svg: '<svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.2" stroke-linecap="round" stroke-linejoin="round"><rect x="3" y="7" width="18" height="13" rx="2"/><path d="M8 7V4h8v3M3 12h18M10 12v2h4v-2"/></svg>'
  },
  security: {
    label: 'จุดรักษาความปลอดภัย', color: '#475569', bg: '#f1f5f9',
    svg: '<svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.3" stroke-linecap="round" stroke-linejoin="round"><path d="M12 22s8-4 8-10V5l-8-3-8 3v7c0 6 8 10 8 10Z"/><path d="m9 12 2 2 4-4"/></svg>'
  },
  firstAid: {
    label: 'จุดปฐมพยาบาล', color: '#dc2626', bg: '#fee2e2',
    svg: '<svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.4" stroke-linecap="round" stroke-linejoin="round"><rect x="3" y="3" width="18" height="18" rx="4"/><path d="M12 8v8M8 12h8"/></svg>'
  },
  fireSafety: {
    label: 'อุปกรณ์ดับเพลิง', color: '#dc2626', bg: '#fee2e2',
    svg: '<svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.3" stroke-linecap="round" stroke-linejoin="round"><path d="M12 22c4 0 7-3 7-7 0-3-2-5-4-7 0 3-2 4-3 4 1-4-1-7-4-10 0 5-3 7-3 12 0 5 3 8 7 8Z"/><path d="M9 18c0-2 1-3 3-5 0 2 2 2 2 4 0 2-1 3-2 3"/></svg>'
  },
  elevator: {
    label: 'ลิฟต์', color: '#7c3aed', bg: '#ede9fe',
    svg: '<svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.2" stroke-linecap="round" stroke-linejoin="round"><rect x="5" y="3" width="14" height="18" rx="2"/><path d="m9 8 3-3 3 3M9 16l3 3 3-3M12 5v14"/></svg>'
  },
  accessible: {
    label: 'ทางลาดผู้ใช้รถเข็น', color: '#2563eb', bg: '#dbeafe',
    svg: '<svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.2" stroke-linecap="round" stroke-linejoin="round"><circle cx="9" cy="5" r="2"/><path d="M8 10h5l2 4h3M9 8v7a5 5 0 1 0 5 5M13 10l2-2"/></svg>'
  }
} as const

const getFacilityMetaForFacility = (facility: MarketFacility) => {
  const name = String(facility.name || '').trim().toLowerCase()
  if (name.includes('ห้องน้ำ') || name.includes('สุขา') || name.includes('ส้วม') || name.includes('toilet') || name.includes('restroom') || /(^|\s)wc($|\s)/i.test(name)) return FACILITY_META.toilet
  if (name.includes('อ่างล้างมือ')) return NAMED_FACILITY_META.sink
  if (name.includes('บันไดเลื่อน')) return NAMED_FACILITY_META.escalator
  if (name.includes('บันได')) return NAMED_FACILITY_META.stairs
  if (name.includes('ลิฟต์') || name.includes('elevator')) return NAMED_FACILITY_META.elevator
  if (name.includes('ทางลาด') || name.includes('รถเข็น') || name.includes('ผู้พิการ') || name.includes('accessible')) return NAMED_FACILITY_META.accessible
  if (name.includes('ที่จอด') || name.includes('parking')) return FACILITY_META.parking
  if (name.includes('ทางเข้า') || name.includes('ทางออก') || name.includes('ประตู') || name.includes('entrance') || name.includes('exit')) return FACILITY_META.exit
  if (name.includes('ขยะ') || name.includes('waste') || name.includes('trash')) return FACILITY_META.waste
  if (name.includes('รปภ') || name.includes('รักษาความปลอดภัย') || name.includes('ป้อมยาม') || name.includes('security')) return NAMED_FACILITY_META.security
  if (name.includes('ปฐมพยาบาล') || name.includes('พยาบาล') || name.includes('first aid')) return NAMED_FACILITY_META.firstAid
  if (name.includes('ดับเพลิง') || name.includes('ถังดับเพลิง') || name.includes('fire')) return NAMED_FACILITY_META.fireSafety
  if (name.includes('ศาลเจ้า') || name.includes('ศาลพระพรหม')) return NAMED_FACILITY_META.shrine
  if (name.includes('หอนาฬิกา')) return NAMED_FACILITY_META.clockTower
  if (name.includes('เต้น') || name.includes('เต็นท์')) return NAMED_FACILITY_META.tent
  if (name.includes('ลาน')) return NAMED_FACILITY_META.plaza
  if (name.includes('ห้องประชุม')) return NAMED_FACILITY_META.meeting
  if (name.includes('ศูนย์ฝึก')) return NAMED_FACILITY_META.training
  if (name.includes('ประชาสัมพันธ์') || name.includes('ประชาสัมพันธุ์') || name.includes('จุดบริการ') || name.includes('information')) return FACILITY_META.service
  if (name.startsWith('งาน') || name.includes('ปฏิบัติการ') || name === 'บก.' || name === 'diy') return NAMED_FACILITY_META.office
  return getFacilityMeta(facility.type)
}

const findHostZoneForFacility = (facility: MarketFacility): ProjectZone | null => {
  if (facility.projectId) {
    const byId = projects.value.find(p => p.id === facility.projectId)
    if (byId) return byId
  }
  // Imported survey facilities without projectId are explicitly market-level.
  // Spatial fallback is reserved for facilities added manually on the map.
  if (facility.sourceProperties) return null
  return projects.value.find(p => pointInProject(facility.coordinates, p)) || null
}

const showPopupForFacility = (facility: MarketFacility) => {
  if (!map || !hoverPopup) return
  if (isAnyEditingToolActive.value) {
    hoverPopup.remove()
    return
  }
  const meta = getFacilityMetaForFacility(facility)
  const hostZone = findHostZoneForFacility(facility)
  const canEditFacility = props.managementMode && !props.inspectOnly
  const deleteAction = canEditFacility
    ? `<button type="button" class="mlp-facility-delete" data-facility-id="${facility.id}"><svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M3 6h18M8 6V4h8v2M19 6l-1 14H6L5 6M10 11v5M14 11v5"/></svg>ลบจุดบริการ</button>`
    : ''
  hoverPopup.setLngLat(facility.coordinates).setHTML(`
    <div class="mlp" style="--mlp-accent:${meta.color}">
      <div class="mlp-body">
        <div class="mlp-head">
          <div class="mlp-title">
            <span style="display:inline-flex;align-items:center;justify-content:center;width:24px;height:24px;border-radius:7px;background:${meta.bg};color:${meta.color}">${meta.svg}</span>
            <strong>${facility.name}</strong>
          </div>
          <span class="mlp-pill" style="background:${meta.bg};color:${meta.color}">${meta.label}</span>
        </div>
        <div class="mlp-sub">${hostZone ? `ตั้งอยู่ใน: ${hostZone.name}` : 'จุดอำนวยความสะดวกส่วนกลางตลาด'}</div>
        <div class="mlp-hint">${canEditFacility ? 'จัดการจุดบริการบนผังตลาด' : 'จุดอำนวยความสะดวกสำหรับผู้ค้าและผู้มาใช้บริการ'}</div>
        ${deleteAction}
      </div>
    </div>
  `).addTo(map)
  const popupElement = hoverPopup.getElement()
  popupElement?.classList.toggle('is-interactive', canEditFacility)
  if (canEditFacility) {
    // Keep the popup open and prevent map handlers from receiving clicks made
    // inside it. Attach to the inner .mlp node so listeners are fresh per-popup.
    const innerMlp = popupElement?.querySelector('.mlp')
    innerMlp?.addEventListener('mousedown', (event: Event) => event.stopPropagation())
    innerMlp?.addEventListener('click', (event: Event) => event.stopPropagation())
    const deleteButton = popupElement?.querySelector(`[data-facility-id="${facility.id}"]`)
    deleteButton?.addEventListener('click', (event: Event) => {
      event.preventDefault()
      event.stopPropagation()
      pinnedFacilityPopup.value = false
      hoverPopup?.remove()
      emit('requestDeleteFacility', facility)
    }, { once: true })
  }
}

const clearFacilityMarkers = () => {
  facilityMarkers.forEach(m => {
    try { m.remove() } catch {}
  })
  facilityMarkers = []
}

const renderFacilityMarkers = () => {
  clearFacilityMarkers()
  if (!map || !maplibregl) return
  if (map.getLayer('facilities-circle')) map.setLayoutProperty('facilities-circle', 'visibility', 'none')
  if (map.getLayer('facilities-label')) map.setLayoutProperty('facilities-label', 'visibility', 'none')
  if (!showFacilities.value) return

  managedFacilities.value
    .filter(facility => {
      const visibleFloorId = activeProjectId.value ? currentFloorId.value : selectedMarketFloorId.value
      return !visibleFloorId
        || visibleFloorId === 'all'
        || !facility.floorId
        || String(facility.floorId) === String(visibleFloorId)
    })
    .forEach(facility => {
      if (!Array.isArray(facility.coordinates) || facility.coordinates.length < 2) return
      const meta = getFacilityMetaForFacility(facility)
      const el = document.createElement('button')
      el.type = 'button'
      el.className = 'facility-map-marker'
      el.setAttribute('aria-label', `${meta.label}: ${facility.name}`)
      el.title = `${meta.label}: ${facility.name}`
      Object.assign(el.style, {
        display: 'inline-flex',
        alignItems: 'center',
        justifyContent: 'center',
        width: '28px',
        height: '28px',
        padding: '0',
        gap: '0',
        borderRadius: '9999px',
        border: `1.5px solid ${meta.color}`,
        background: '#ffffff',
        color: '#0f172a',
        fontFamily: "'Bai Jamjuree', sans-serif",
        boxShadow: '0 2px 6px rgba(15,23,42,0.16)',
        cursor: 'pointer',
        whiteSpace: 'nowrap',
        zIndex: '8'
      })
      el.innerHTML = `
        <span style="display:inline-flex;align-items:center;justify-content:center;flex:0 0 20px;width:20px;height:20px;border-radius:50%;background:${meta.bg};color:${meta.color}">${meta.svg}</span>
      `
      const markerSvg = el.querySelector('svg')
      if (markerSvg) {
        markerSvg.style.width = '16px'
        markerSvg.style.height = '16px'
      }
      el.addEventListener('mouseenter', () => {
        showPopupForFacility(facility)
      })
      el.addEventListener('mouseleave', () => {
        if (!props.managementMode) hoverPopup?.remove()
      })
      el.addEventListener('click', (ev) => {
        ev.stopPropagation()
        pinnedFacilityPopup.value = props.managementMode
        showPopupForFacility(facility)
      })
      const marker = new maplibregl.Marker({ element: el, anchor: 'center' })
        .setLngLat(facility.coordinates)
        .addTo(map)
      facilityMarkers.push(marker)
    })
}

// ---------------- Popups ----------------
const statusMeta: Record<string, { label: string; bg: string; fg: string }> = {
  available: { label: 'ว่าง', bg: '#dcfce7', fg: '#15803d' },
  reserved: { label: 'ติดจอง', bg: '#fef3c7', fg: '#d97706' },
  occupied: { label: 'มีผู้เช่าแล้ว', bg: '#f1f5f9', fg: '#64748b' }
}

const ICONS = {
  area: '<svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><rect x="3" y="3" width="18" height="18" rx="2"></rect><path d="M3 9h18M9 21V9"></path></svg>',
  grid: '<svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><rect x="3" y="3" width="7" height="7"></rect><rect x="14" y="3" width="7" height="7"></rect><rect x="14" y="14" width="7" height="7"></rect><rect x="3" y="14" width="7" height="7"></rect></svg>',
  ruler: '<svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M21.3 15.3a2.4 2.4 0 0 1 0 3.4l-2.6 2.6a2.4 2.4 0 0 1-3.4 0L2.7 8.7a2.4 2.4 0 0 1 0-3.4l2.6-2.6a2.4 2.4 0 0 1 3.4 0Z"></path><path d="m14.5 12.5 2-2M11.5 9.5l2-2M8.5 6.5l2-2M17.5 15.5l2-2"></path></svg>',
  tag: '<svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M12 2H2v10l9.29 9.29a1 1 0 0 0 1.41 0l8.59-8.59a1 1 0 0 0 0-1.41Z"></path><circle cx="7" cy="7" r="1.5"></circle></svg>',
  wallet: '<svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M19 7V4a1 1 0 0 0-1-1H5a2 2 0 0 0 0 4h15a1 1 0 0 1 1 1v4h-3a2 2 0 0 0 0 4h3a1 1 0 0 0 1-1v-2"></path><path d="M3 5v14a2 2 0 0 0 2 2h15a1 1 0 0 0 1-1v-4"></path></svg>',
  key: '<svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="7.5" cy="15.5" r="5.5"></circle><path d="m21 2-9.6 9.6M15.5 7.5l3 3L22 7l-3-3"></path></svg>',
  arrow: '<svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.2" stroke-linecap="round" stroke-linejoin="round"><path d="M5 12h14M12 5l7 7-7 7"></path></svg>',
  lock: '<svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><rect x="3" y="11" width="18" height="11" rx="2"></rect><path d="M7 11V7a5 5 0 0 1 10 0v4"></path></svg>',
  clock: '<svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="10"></circle><polyline points="12 6 12 12 16 14"></polyline></svg>'
}

const showPopupForProject = (p: ProjectZone, lngLat: [number, number]) => {
  if (!map || !hoverPopup) return
  if (isAnyEditingToolActive.value) {
    hoverPopup.remove()
    return
  }
  const layout = getProjectLayout(p)
  const total = layout?.total ?? 0
  const available = layout?.available ?? 0
  const reserved = layout?.reserved ?? 0
  const occupied = Math.max(total - available - reserved, 0)
  const pct = (n: number) => (total ? (n / total) * 100 : 0)
  const isFull = available === 0
  const accent = p.color

  const hasFloors = p.floors && p.floors.length > 1
  const floorBadge = hasFloors
    ? `<div class="mlp-building-row"><span class="mlp-building-badge">อาคาร ${p.floors.length} ชั้น</span></div>`
    : ''
  const hintText = hasFloors
    ? 'อาคารนี้มี ' + p.floors.length + ' ชั้น: คลิกเพื่อเลือกระดับชั้นก่อนเลือกแผงค้า'
    : (isFull ? 'คลิกที่กรอบโครงการเพื่อดูแผงย่อยและเข้าคิวสำรอง' : 'คลิกที่กรอบโครงการเพื่อดูแผงย่อยและเลือกจอง')

  hoverPopup.setLngLat(lngLat).setHTML(`
    <div class="mlp" style="--mlp-accent:${accent}">
      <div class="mlp-body">
        <div class="mlp-head">
          <div class="mlp-title">
            <span class="mlp-dot"></span>
            <strong>${p.name}</strong>
          </div>
          <span class="mlp-pill ${isFull ? 'is-full' : 'is-open'}">${isFull ? 'เต็มแล้ว' : `ว่าง ${available} แผง`}</span>
        </div>
        ${floorBadge}
        <div class="mlp-sub">${p.category}</div>

        <div class="mlp-tiles">
          <div class="mlp-tile">
            <span class="mlp-tile-label">${ICONS.area}พื้นที่โครงการ</span>
            <b class="mlp-tile-value is-primary">${p.realAreaSqm.toLocaleString()} <small>ตร.ม.</small></b>
          </div>
          <div class="mlp-tile">
            <span class="mlp-tile-label">${ICONS.grid}แผงย่อยทั้งหมด</span>
            <b class="mlp-tile-value">${total.toLocaleString()} <small>แผง</small></b>
          </div>
        </div>

        <div class="mlp-occ">
          <div class="mlp-occ-head">
            <span>สถานะแผงในโครงการ</span>
            <span class="mlp-occ-pct">เช่าแล้ว ${Math.round(pct(occupied))}%</span>
          </div>
          <div class="mlp-occ-bar">
            <i style="width:${pct(available)}%;background:#10b981"></i>
            <i style="width:${pct(reserved)}%;background:#f59e0b"></i>
            <i style="width:${pct(occupied)}%;background:#475569"></i>
          </div>
          <div class="mlp-occ-legend">
            <span><i style="background:#10b981"></i>ว่าง <b style="color:#059669">${available}</b></span>
            <span><i style="background:#f59e0b"></i>ติดจอง <b style="color:#d97706">${reserved}</b></span>
            <span><i style="background:#475569"></i>มีผู้เช่า <b>${occupied}</b></span>
          </div>
        </div>

        <div class="mlp-hint">${hintText}</div>
      </div>
    </div>
  `).addTo(map)
  hoverPopup.getElement()?.classList.remove('is-interactive')
}

const stallPhotoSamples = [
  '/images/stall_3d_perspective.jpg',
  '/images/market_stall_sample.jpg',
  '/images/market_hall_front.jpg'
]

const getStallDisplayPhoto = (s: MapStall): string => {
  if (s.photo && s.photo !== '/images/stall_3d_perspective.jpg') return s.photo
  let hash = 0
  for (let i = 0; i < s.code.length; i++) hash = ((hash << 5) - hash) + s.code.charCodeAt(i)
  return stallPhotoSamples[Math.abs(hash) % stallPhotoSamples.length]
}

const showPopupForStall = (s: MapStall) => {
  if (!map || !hoverPopup || !s?.center) return
  const displayName = getStallDisplayName(s)
  if (!displayName) {
    hoverPopup.remove()
    return
  }
  hoverPopup.setLngLat(s.center).setHTML(`
    <div class="mlp mlp-stall-tooltip"><strong>แผง ${displayName}</strong></div>
  `).addTo(map)
}

// ---------------- Drill-down actions ----------------
let selectedStallMarkers: any[] = []

const renderSelectedStallPins = () => {
  if (!map || !maplibregl) return
  selectedStallMarkers.forEach(m => m.remove())
  selectedStallMarkers = []

  // The selected-stall fill/stroke is sufficient; do not cover its location
  // with a floating check marker. Keep a beacon only when viewing another project.
  if (props.previewMode && bookedProject.value && activeProjectId.value !== bookedProjectId.value) {
    // หากอยู่ใน Step 3 แต่กำลังคลิกดูโครงการอื่น ให้ปักหมุด Beacon นำทางบนโครงการที่มีแผงที่จองไว้
    const bookedProject = projects.value.find(project => project.id === bookedProjectId.value)
    const bLayout = bookedProject ? getProjectLayout(bookedProject) : null
    const bookedStallsList = (bLayout?.stalls || []).filter(isStallSelected)
    const pinLocation: [number, number] = bookedStallsList.length > 0
      ? bookedStallsList[0].center
      : bookedProject.value.center

    const el = document.createElement('div')
    el.className = 'booked-project-beacon-marker'
    el.innerHTML = `
      <div class="beacon-pin-bubble" title="คลิกเพื่อกลับไปยังแผงที่คุณเลือก (${bookedProject.value.name})">
        <span class="beacon-pin-icon">
          <span class="beacon-pulse-ring"></span>
          <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="#ffffff" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round">
            <path d="M21 10c0 7-9 13-9 13s-9-6-9-13a9 9 0 0 1 18 0z"></path>
            <circle cx="12" cy="10" r="3"></circle>
          </svg>
        </span>
        <span>แผงที่คุณจอง (${props.selectedStallCodes.join(', ')})</span>
      </div>
    `
    el.onclick = () => focusBookedStalls()

    const marker = new maplibregl.Marker({ element: el, anchor: 'bottom' })
      .setLngLat(pinLocation)
      .addTo(map)
    selectedStallMarkers.push(marker)
  }
}

const focusSelectedStalls = (): boolean => {
  const targetProject = bookedProject.value
  if (!targetProject) return false
  let targetFloorId: string | undefined
  if (targetProject.floors && targetProject.floors.length > 0 && props.selectedStallCodes?.length) {
    for (const fl of targetProject.floors) {
      const flLayout = applyMarketStatus(getProjectStallLayout(targetProject, fl.id))
      if (flLayout.stalls.some(isStallSelected)) {
        targetFloorId = fl.id
        break
      }
    }
  }
  enterProject(targetProject, targetFloorId)
  return true
}

const focusBookedStalls = focusSelectedStalls

const getDrillMapPadding = () => {
  const isResponsive = typeof window !== 'undefined' && window.innerWidth <= 768
  if (props.previewMode) {
    return isResponsive
      ? { top: 50, bottom: 50, left: 20, right: 20 }
      : { top: 60, bottom: 60, left: 60, right: 60 }
  }
  if (isResponsive) {
    const sheetH = projectSheetHeight.value || (typeof window !== 'undefined' ? window.innerHeight * 0.5 : 360)
    const bottomPad = Math.round(Math.min(window.innerHeight * 0.56, Math.max(240, sheetH + 16)))
    return { top: 70, bottom: bottomPad, left: 20, right: 20 }
  }
  // ลดพื้นที่กันขอบลงเล็กน้อย เพื่อให้ตอนเข้าโครงการเห็นแผงใกล้ขึ้น
  // โดยยังเหลือพื้นที่สำหรับการ์ดซ้ายและแผงเครื่องมือด้านขวา
  if (appliedSearchQuery.value.trim()) return { top: 55, bottom: 35, left: 40, right: 340 }
  return { top: 55, bottom: 35, left: 280, right: 340 }
}

const recenterMapView = () => {
  if (!map) return
  hoverPopup?.remove()
  if (activeProject.value) {
    const p = activeProject.value
    const pLayout = getProjectLayout(p) || activeLayout.value
    const selectedStallsList = (pLayout?.stalls || []).filter(isStallSelected)
    const stallBounds = props.previewMode && selectedStallsList.length > 0 ? getStallsBounds(selectedStallsList) : null

    map.fitBounds(stallBounds || getProjectBounds(p), {
      padding: getDrillMapPadding(),
      maxZoom: 20.6,
      pitch: DRILL_PITCH,
      bearing: currentMarketBearing.value,
      duration: 800
    })
  } else {
    map.flyTo({
      center: currentMarketCenter.value,
      zoom: currentMarketZoom.value,
      pitch: currentMarketPitch.value,
      bearing: currentMarketBearing.value,
      duration: 800
    })
  }
}

const openProject = (p: ProjectZone) => {
  if (isDrawingZone.value || isDrawingStall.value || isAddingFacility.value) return
  if (props.managementMode && isPolygonEditMode.value && !activeProjectId.value) {
    overviewEditProjectId.value = p.id
    hoverPopup?.remove()
    renderPolygonVertexHandles()
    refreshSources()
    return
  }
  if (!props.previewMode && !props.managementMode && (p.floors?.length || 0) > 1) {
    hoverPopup?.remove()
    floorSelectorProject.value = p
    return
  }
  enterProject(p)
}

const closeFloorSelector = () => {
  floorSelectorProject.value = null
}

const selectFloorFromDialog = (floorId: string) => {
  const project = floorSelectorProject.value
  if (!project) return
  floorSelectorProject.value = null
  enterProject(project, floorId)
}

const enterProject = (p: ProjectZone, directToFloor?: string) => {
  floorSelectorProject.value = null
  hoverPopup?.remove()
  if (!activeProjectId.value) marketFloorBeforeProject.value = selectedMarketFloorId.value
  const projectFloorIds = getProjectActualFloorIds(p)
  if (marketFloors.value.length > 1 && !projectFloorIds.includes(selectedMarketFloorId.value)) {
    selectedMarketFloorId.value = directToFloor && projectFloorIds.includes(directToFloor)
      ? directToFloor
      : (projectFloorIds[0] || marketFloors.value[0]?.id || '')
  }
  activeProjectId.value = p.id
  overviewEditProjectId.value = p.id
  if (typeof window !== 'undefined' && window.innerWidth <= 768) {
    ovPanel.value = null
    // ทุกครั้งที่เปิดโครงการบนมือถือ ให้เริ่มเป็น Bottom Sheet ครึ่งหน้าจอเสมอ
    // ป้องกันค่าความสูงจากการลากหรือขยายโครงการก่อนหน้าค้างอยู่
    isProjectSheetExpanded.value = false
    projectSheetHeight.value = Math.round(window.innerHeight * 0.5)
  }
  isPolygonEditMode.value = false
  isStallPolygonEditMode.value = false
  editingStallCode.value = null
  clearPolygonVertexHandles()
  clearStallVertexHandles()
  hoveredProject.value = null
  hoveredStallKey.value = null

  if (directToFloor) {
    currentFloorId.value = directToFloor
  } else if (!currentFloorId.value || !p.floors?.some(f => f.id === currentFloorId.value)) {
    // Smart auto-select: เลือกชั้นแรกที่มีแผงว่างให้อัตโนมัติ (ไม่เลือกชั้นที่เต็ม 0 แผงเป็นชั้นแรก)
    const firstAvailFloor = p.floors?.find(f => getFloorLayout(p, f.id).available > 0)
    currentFloorId.value = firstAvailFloor?.id || p.floors?.[0]?.id || '1'
  }

  // Show the selected project's stall polygons immediately. DOM facility
  // markers remain deferred until the camera settles because rebuilding those
  // markers is the expensive part that can interrupt the fit animation.
  const deferPublicProjectDetails = !props.managementMode && !props.previewMode
  if (deferPublicProjectDetails && map) {
    map.getSource('projects-polygons')?.setData(getProjectsGeoJson())
    map.getSource('project-stalls')?.setData(getStallsGeoJson())
  } else {
    refreshSources()
  }
  setBuildingsOpacity(BUILDINGS_OPACITY_DRILL)

  const pLayout = getProjectLayout(p) || activeLayout.value
  const selectedStallsList = (pLayout?.stalls || []).filter(isStallSelected)
  const stallBounds = props.previewMode && selectedStallsList.length > 0 ? getStallsBounds(selectedStallsList) : null

  if (deferPublicProjectDetails && map) {
    const expectedProjectId = p.id
    map.once('moveend', () => {
      if (isUnmounted || activeProjectId.value !== expectedProjectId) return
      map.getSource('market-facilities')?.setData(getFacilitiesGeoJson())
      renderFacilityMarkers()
      renderSelectedStallPins()
    })
  }

  map?.fitBounds(stallBounds || getProjectBounds(p), {
    padding: getDrillMapPadding(),
    maxZoom: 20.6,
    pitch: DRILL_PITCH,
    bearing: map.getBearing(),
    duration: 800
  })
  if (!deferPublicProjectDetails) renderSelectedStallPins()
  renderPolygonVertexHandles()
  // ส่งสถิติจาก layout ที่ผ่านการรวมสถานะจริงแล้วให้หน้าจองด้วย
  // ป้องกัน empty state ใช้ค่า stallsAvailable เดิมจาก master project
  // ในขณะที่การ์ดสรุปด้านซ้ายแสดงว่าพื้นที่เต็มแล้ว
  emit('selectProject', {
    ...p,
    totalStalls: pLayout.total,
    stallsAvailable: pLayout.available
  })
}

// ควบคุมการเลื่อนปรับขนาดความสูง Bottom Sheet ของโครงการบนโหมด Responsive
const isProjectSheetExpanded = ref(false)
const projectSheetHeight = ref<number | null>(null)
const isSheetDragging = ref(false)
let sheetTouchStartY = 0
let sheetTouchStartHeight = 0
let sheetMoved = false

const toggleProjectSheetExpand = () => {
  isProjectSheetExpanded.value = !isProjectSheetExpanded.value
  projectSheetHeight.value = isProjectSheetExpanded.value ? Math.round(window.innerHeight - 58) : Math.round(window.innerHeight * 0.5)
}

const onSheetTouchStart = (e: TouchEvent) => {
  if (typeof window === 'undefined' || window.innerWidth > 768) return
  if (e.touches && e.touches.length > 0) {
    sheetTouchStartY = e.touches[0].clientY
    sheetMoved = false
    isSheetDragging.value = true
    const panel = document.querySelector('.stall-mode-panel') as HTMLElement | null
    sheetTouchStartHeight = panel ? panel.offsetHeight : (projectSheetHeight.value || window.innerHeight * 0.5)
  }
}

const onSheetTouchMove = (e: TouchEvent) => {
  if (!isSheetDragging.value || typeof window === 'undefined' || window.innerWidth > 768) return
  if (e.touches && e.touches.length > 0) {
    const currentY = e.touches[0].clientY
    const deltaY = sheetTouchStartY - currentY
    if (Math.abs(deltaY) > 4) {
      sheetMoved = true
    }
    const calculatedHeight = sheetTouchStartHeight + deltaY
    const minHeight = window.innerHeight * 0.18
    const maxHeight = window.innerHeight - 58
    projectSheetHeight.value = Math.min(maxHeight, Math.max(minHeight, calculatedHeight))
  }
}

const onSheetTouchEnd = () => {
  if (!isSheetDragging.value || typeof window === 'undefined' || window.innerWidth > 768) return
  isSheetDragging.value = false
  if (!sheetMoved) {
    // ถ้าแค่แตะคลิกเบาๆ โดยไม่ได้ลาก: สลับระหว่างครึ่งจอกับเต็มจอใต้ Navbar
    const currentH = projectSheetHeight.value || window.innerHeight * 0.5
    if (currentH < window.innerHeight * 0.75) {
      projectSheetHeight.value = Math.round(window.innerHeight - 58)
      isProjectSheetExpanded.value = true
    } else {
      projectSheetHeight.value = Math.round(window.innerHeight * 0.5)
      isProjectSheetExpanded.value = false
    }
    return
  }

  // หากลากลงมาต่ำกว่า 22vh หรือรูดลงมาใกล้ขอบจอล่างสุด -> ให้ปิดกลับผังรวมตลาด
  if (projectSheetHeight.value && projectSheetHeight.value <= window.innerHeight * 0.22) {
    projectSheetHeight.value = 0
    setTimeout(() => {
      projectSheetHeight.value = null
      isProjectSheetExpanded.value = false
      exitProject()
    }, 180)
  } else {
    // ผู้ใช้เลื่อนปรับขนาดไว้ตรงไหน ก็ให้ค้างขนาดไว้ตามที่ผู้ใช้ปรับเลย
    if (projectSheetHeight.value && projectSheetHeight.value > window.innerHeight * 0.7) {
      isProjectSheetExpanded.value = true
    } else {
      isProjectSheetExpanded.value = false
    }
  }
}

const handleBackOverviewClick = () => {
  if (isProjectPanelExiting.value) return
  isProjectPanelExiting.value = true
  window.setTimeout(() => {
    isProjectPanelExiting.value = false
    exitProject()
  }, 180)
}

const exitProject = () => {
  isProjectSheetExpanded.value = false
  projectSheetHeight.value = null
  isSheetDragging.value = false
  finishActiveTools()
  floorSelectorProject.value = null
  hoverPopup?.remove()
  activeProjectId.value = null
  overviewEditProjectId.value = null
  isPolygonEditMode.value = false
  isStallPolygonEditMode.value = false
  editingStallCode.value = null
  clearPolygonVertexHandles()
  clearStallVertexHandles()
  hoveredStallKey.value = null
  currentFloorId.value = '1'
  selectedMarketFloorId.value = marketFloorBeforeProject.value
  refreshSources()
  setBuildingsOpacity(BUILDINGS_OPACITY)
  map?.flyTo({
    center: currentMarketCenter.value,
    zoom: currentMarketZoom.value,
    pitch: currentMarketPitch.value,
    bearing: currentMarketBearing.value,
    duration: 700
  })
  emit('exitProject')
}

// Called from the info banner button
const selectProject = (p: ProjectZone) => openProject(p)

const findStall = (code: string): MapStall | null => {
  if (!activeProject.value) return null
  const inCurrent = activeLayout.value?.stalls.find(s => s.code === code)
  if (inCurrent) return inCurrent
  if (activeProject.value.floors && activeProject.value.floors.length > 1) {
    for (const fl of activeProject.value.floors) {
      const flLayout = applyMarketStatus(getProjectStallLayout(activeProject.value, fl.id))
      const found = flLayout.stalls.find(s => s.code === code)
      if (found) return found
    }
  }
  return null
}

const getStallFromFeature = (feature: any): MapStall | null => {
  const layout = activeLayout.value
  const stallIndex = Number(feature?.properties?.stallIndex)

  // `code` is not a reliable map key: imported GIS data can contain repeated
  // labels. `stallIndex` is generated from this exact GeoJSON feature.
  if (layout && Number.isInteger(stallIndex) && stallIndex >= 0) {
    const stall = layout.stalls[stallIndex]
    if (stall) return stall
  }

  const code = feature?.properties?.code
  return typeof code === 'string' ? findStall(code) : null
}

const setHoveredProject = (project: ProjectZone | null): boolean => {
  if (hoveredProject.value?.id === project?.id) return false
  projectMarkerElements.forEach((element, projectId) => {
    element.classList.toggle('is-project-hovered', projectId === project?.id)
  })
  hoveredProject.value = project
  // Update just the project source so the polygon outline follows the cursor.
  map?.getSource('projects-polygons')?.setData(getProjectsGeoJson())
  return true
}

const getMostSpecificProjectHit = (hits: any[]): ProjectZone | null => {
  const matchingProjects = hits
    .map(hit => projects.value.find(project => project.id === hit.properties.id))
    .filter((project): project is ProjectZone => Boolean(project))

  if (!matchingProjects.length) return null
  // MapLibre returns features in paint order. The smallest footprint is the
  // most specific zone when surveyed polygons overlap.
  return matchingProjects.reduce((smallest, project) => (
    projectFootprintArea(project) < projectFootprintArea(smallest) ? project : smallest
  ))
}

const getFacilityFromFeature = (feature: any): MarketFacility | null => {
  const facilityId = feature?.properties?.id
  if (typeof facilityId !== 'string') return null
  return managedFacilities.value.find(facility => facility.id === facilityId) || null
}

const bindMapEvents = () => {
  map.on('zoomend', () => {
    // Labels are adaptive: redraw only when the map settles after zooming.
    renderZoneLabels()
  })

  map.on('mousemove', (e: any) => {
    // A facility popup selected by a staff member is an action surface, not a
    // hover hint. Do not let the regular map hover logic dismiss it.
    if (pinnedFacilityPopup.value) return
    if (isDrawingZone.value || isDrawingStall.value || isAddingFacility.value) {
      drawingCursor.value = [e.lngLat.lng, e.lngLat.lat]
      refreshDrawingZone()
      map.getCanvas().style.cursor = 'crosshair'
      return
    }
    const facilityHit = map.getLayer('facilities-polygon-fill')
      ? map.queryRenderedFeatures(e.point, { layers: ['facilities-polygon-fill'] })
      : []
    if (facilityHit.length) {
      setHoveredProject(null)
      setHoveredStall(null)
      const facility = getFacilityFromFeature(facilityHit[0])
      map.getCanvas().style.cursor = facility ? 'pointer' : ''
      if (facility) showPopupForFacility(facility)
      return
    }
    const stallHit = map.getLayer('stalls-fill')
      ? map.queryRenderedFeatures(e.point, { layers: ['stalls-fill'] })
      : []
    if (stallHit.length) {
      setHoveredProject(null)
      const s = getStallFromFeature(stallHit[0])
      const changedStall = setHoveredStall(s)
      const isMyBooked = Boolean(s && !props.previewMode && (props.myBookedStallCodes || []).includes(s.code))
      const canClick = Boolean(s && (props.managementMode || (s.status !== 'occupied' && !isMyBooked)))
      map.getCanvas().style.cursor = canClick ? 'pointer' : 'not-allowed'
      // Keep stall details click-driven so the tooltip does not flicker while moving across the map.
      return
    }

    setHoveredStall(null)

    const projHit = map.getLayer('projects-fill')
      ? map.queryRenderedFeatures(e.point, { layers: ['projects-fill'] })
      : []
    const p = getMostSpecificProjectHit(projHit)
    if (p && p.id !== activeProjectId.value) {
      map.getCanvas().style.cursor = 'pointer'
      const changedProject = setHoveredProject(p)
      if (changedProject) showPopupForProject(p, [e.lngLat.lng, e.lngLat.lat])
    } else {
      map.getCanvas().style.cursor = ''
      setHoveredProject(null)
      hoverPopup.remove()
    }
  })

  map.on('mouseout', () => {
    if (pinnedFacilityPopup.value) return
    setHoveredProject(null)
    setHoveredStall(null)
    hoverPopup?.remove()
  })

  map.on('click', (e: any) => {
    if (pinnedFacilityPopup.value) {
      pinnedFacilityPopup.value = false
      hoverPopup?.remove()
      return
    }
    if (isAddingFacility.value) {
      const coordinates: [number, number] = [e.lngLat.lng, e.lngLat.lat]
      const projHitForFacility = map.getLayer('projects-fill')
        ? map.queryRenderedFeatures(e.point, { layers: ['projects-fill'] })
        : []
      const hostProject = activeProject.value || getMostSpecificProjectHit(projHitForFacility)
      cancelDrawing()
      emit('requestAddFacility', coordinates, hostProject?.id)
      return
    }
    if (isDrawingZone.value || isDrawingStall.value) {
      if (isDrawingStall.value && stallDrawShape.value === 'rectangle') {
        const clickedPt: [number, number] = [e.lngLat.lng, e.lngLat.lat]
        if (!rectStartCorner.value || drawingPoints.value.length >= 4) {
          rectStartCorner.value = clickedPt
          drawnStallDimensions.value = null
          drawingPoints.value = [clickedPt]
          refreshDrawingZone()
          return
        }
        const { corners, widthM, depthM } = buildRectangleCorners(rectStartCorner.value, clickedPt)
        drawingPoints.value = corners
        drawnStallDimensions.value = { widthM, depthM }
        refreshDrawingZone()
        return
      }

      if (drawingPoints.value.length >= 3 && drawingPoints.value[0]) {
        const firstScreenPt = map.project(drawingPoints.value[0])
        const dist = Math.hypot(e.point.x - firstScreenPt.x, e.point.y - firstScreenPt.y)
        if (dist <= 16) {
          finishDrawing()
          return
        }
      }
      drawingPoints.value.push([e.lngLat.lng, e.lngLat.lat])
      refreshDrawingZone()
      return
    }

    const facilityHit = map.getLayer('facilities-polygon-fill')
      ? map.queryRenderedFeatures(e.point, { layers: ['facilities-polygon-fill'] })
      : []
    if (facilityHit.length) {
      const facility = getFacilityFromFeature(facilityHit[0])
      if (facility) {
        pinnedFacilityPopup.value = props.managementMode
        showPopupForFacility(facility)
      }
      return
    }

    const stallHit = map.getLayer('stalls-fill')
      ? map.queryRenderedFeatures(e.point, { layers: ['stalls-fill'] })
      : []
    if (stallHit.length) {
      const s = getStallFromFeature(stallHit[0])
      if (props.managementMode && isDeletingStall.value && s) {
        triggerDeleteStall(s)
        return
      }
      if (props.managementMode && isStallPolygonEditMode.value && s) {
        editingStallCode.value = s.code
        renderStallVertexHandles()
        return
      }
      if (props.previewMode) {
        if (s) showPopupForStall(s)
        return
      }
      if (s && (s.status !== 'occupied' || props.managementMode)) {
        emit('selectStall', s)
      }
      return
    }

    const projHit = map.getLayer('projects-fill')
      ? map.queryRenderedFeatures(e.point, { layers: ['projects-fill'] })
      : []
    if (projHit.length) {
      const p = getMostSpecificProjectHit(projHit)
      if (p && p.id !== activeProjectId.value) {
        openProject(p)
        return
      }
    }
    // คลิกพื้นที่ว่าง (หรือพื้นที่โครงการที่เปิดอยู่) บนหน้าข้อมูลตลาด = ล้างการเลือกแผง กลับไปแผงโครงการ
    if (props.isOverviewPage && activeProjectId.value) emit('clearStallSelection')
  })
}

// ---------------- Area badges inside each frame ----------------
const renderZoneLabels = () => {
  if (!map || !maplibregl) return
  projectMarkers.forEach(m => m.remove())
  projectMarkers = []
  projectMarkerElements.clear()
  if (!showProjectLabels.value) return

  const zoom = map.getZoom()
  const visibleProjects = overviewProjects.value.filter(p => p.id !== activeProjectId.value && !isProjectDimmed(p))
  const useCompactLabels = zoom < 18.2

  // At broad zoom, use only small project numbers. They preserve the map context
  // while making it explicit that each coloured frame is a selectable project.
  if (zoom < 15.5) return

  visibleProjects.forEach((p) => {
    const hasProjectNumber = Boolean(String(p.sectionNumber || '').trim())
    const showFloorContext = !hasProjectNumber || (marketFloors.value.length > 1 && selectedMarketFloorId.value === 'all')
    const projectFloorLabel = getProjectActualFloorIds(p)
      .map(floorId => {
        const fl = marketFloors.value.find(floor => floor.id === floorId)
        return getFloorDisplayLabel(fl || floorId)
      })
      .join('/')
    const isCompact = useCompactLabels

    const el = document.createElement('div')
    el.className = 'maplibre-zone-number-badge'
    if (isCompact) el.classList.add('is-compact')
    if (showFloorContext) el.classList.add('has-floor-context')
    if (!hasProjectNumber) el.classList.add('is-floor-only')
    if (hoveredProject.value?.id === p.id) el.classList.add('is-project-hovered')
    el.title = p.name
    projectMarkerElements.set(p.id, el)

    el.innerHTML = `
      <div class="zone-badge-inner">
        ${hasProjectNumber ? `<span class="badge-num">${isCompact ? p.sectionNumber : `โครงการที่ ${p.sectionNumber}`}</span>` : ''}
        ${showFloorContext ? `<span class="badge-floor">${projectFloorLabel}</span>` : ''}
      </div>
    `
    el.addEventListener('click', (ev) => {
      ev.stopPropagation()
      openProject(p)
    })
    el.addEventListener('mouseenter', () => {
      setHoveredProject(p)
      map.getCanvas().style.cursor = 'pointer'
      const lngLat = projectToLngLat(p)
      showPopupForProject(p, lngLat)
    })
    el.addEventListener('mouseleave', () => {
      if (hoveredProject.value?.id === p.id) setHoveredProject(null)
      map.getCanvas().style.cursor = ''
      hoverPopup?.remove()
    })

    // ตรึง Marker ให้อยู่ที่จุดศูนย์กลางของ Polygon ตัวเองเสมอ (offset [0, 0])
    const marker = new maplibregl.Marker({ element: el, anchor: 'center' })
      .setLngLat(projectToLngLat(p))
      .addTo(map)

    projectMarkers.push(marker)
  })
}

const matchesSearchWith = (p: ProjectZone, value: string): boolean => {
  const query = value.trim().toLowerCase()
  if (!query) return true
  if (p.name.toLowerCase().includes(query)) return true
  if (p.category.toLowerCase().includes(query)) return true
  if (props.managementMode && Array.from(getEffectiveFloorLayoutsForProject(p).values())
    .some(layout => layout.stalls.some(stall => stallMatchesSearchWith(stall, query)))) return true
  const numInQuery = query.replace(/[^0-9]/g, '')
  const numInProject = p.name.replace(/[^0-9]/g, '')
  if (numInQuery && numInProject && numInProject.includes(numInQuery)) return true
  return false
}

const matchesSearch = (p: ProjectZone): boolean => matchesSearchWith(p, appliedSearchQuery.value)

const handleSearchInput = () => {
  // Search text is a draft. Keep the map unchanged until the user chooses a suggestion.
  focusedSuggestionIndex.value = -1
  selectedSuggestionId.value = null
  if (searchDebounceTimer) window.clearTimeout(searchDebounceTimer)
  searchDebounceTimer = window.setTimeout(() => {
    suggestionQuery.value = searchQuery.value
    searchDebounceTimer = null
  }, 250)
}

const returnToProjectOverview = () => {
  if (activeProjectId.value) {
    exitProject()
    return
  }
  map?.flyTo({
    center: currentMarketCenter.value,
    zoom: currentMarketZoom.value,
    pitch: currentMarketPitch.value,
    bearing: currentMarketBearing.value,
    duration: 700
  })
}

const clearSearch = () => {
  if (searchDebounceTimer) window.clearTimeout(searchDebounceTimer)
  searchDebounceTimer = null
  searchQuery.value = ''
  suggestionQuery.value = ''
  focusedSuggestionIndex.value = -1
  selectedSuggestionId.value = null
  appliedSearchQuery.value = ''
  refreshSources()
  returnToProjectOverview()
}

const clearAppliedSearch = () => {
  appliedSearchQuery.value = ''
  selectedSuggestionId.value = null
  refreshSources()
  returnToProjectOverview()
}

const clearAllFilters = () => {
  if (searchDebounceTimer) window.clearTimeout(searchDebounceTimer)
  searchDebounceTimer = null
  searchQuery.value = ''
  suggestionQuery.value = ''
  focusedSuggestionIndex.value = -1
  selectedSuggestionId.value = null
  rentFilter.value = 'all'
  sizeFilter.value = 'all'
  availabilityFilter.value = props.stallStatusFilter
  zoneFilter.value = 'all'
  readinessFilter.value = 'all'
  floorFilter.value = 'all'
  contractStatusFilter.value = 'all'
  contractExpiryFilter.value = 'all'
  contractExpiryFrom.value = ''
  contractExpiryTo.value = ''
  bookingRangeFilter.value = 'all'
  bookingFrom.value = ''
  bookingTo.value = ''
  renewalStatusFilter.value = 'all'
  returnStatusFilter.value = 'all'
  legalStatusFilter.value = 'all'
  maintenanceStatusFilter.value = 'all'
  requestTypeFilter.value = 'all'
  draftCategoryFilter.value = ''
  appliedSearchQuery.value = ''
  appliedRentFilter.value = 'all'
  appliedSizeFilter.value = 'all'
  appliedAvailabilityFilter.value = props.stallStatusFilter
  appliedZoneFilter.value = 'all'
  appliedReadinessFilter.value = 'all'
  appliedFloorFilter.value = 'all'
  appliedContractStatusFilter.value = 'all'
  appliedContractExpiryFilter.value = 'all'
  appliedContractExpiryFrom.value = ''
  appliedContractExpiryTo.value = ''
  appliedBookingRangeFilter.value = 'all'
  appliedBookingFrom.value = ''
  appliedBookingTo.value = ''
  appliedRenewalStatusFilter.value = 'all'
  appliedReturnStatusFilter.value = 'all'
  appliedLegalStatusFilter.value = 'all'
  appliedMaintenanceStatusFilter.value = 'all'
  appliedRequestTypeFilter.value = 'all'
  appliedCategoryFilter.value = ''
  emit('update:categoryFilter', null)
  refreshSources()
  returnToProjectOverview()
}

const resetCostFiltersAndOverview = () => {
  emit('reset-cost-filters')
  returnToProjectOverview()
}

const focusAppliedProjectFilterResults = (): boolean => {
  if (!hasActiveFilters.value || appliedSearchQuery.value.trim()) return false
  const matchingProjects = overviewProjects.value.filter(projectMatchesAllFilters)
  if (activeProjectId.value) exitProject()
  if (matchingProjects.length !== 1 || !matchingProjects[0]) {
    returnToProjectOverview()
    return true
  }
  const target = matchingProjects[0]
  nextTick(() => {
    map?.fitBounds(getProjectBounds(target), {
      padding: { top: 80, bottom: 80, left: 80, right: 380 },
      maxZoom: 17.8,
      pitch: currentMarketPitch.value,
      bearing: currentMarketBearing.value,
      duration: 800
    })
  })
  return true
}

const applyFilters = () => {
  appliedRentFilter.value = rentFilter.value
  appliedSizeFilter.value = sizeFilter.value
  appliedAvailabilityFilter.value = availabilityFilter.value
  appliedZoneFilter.value = zoneFilter.value
  appliedReadinessFilter.value = readinessFilter.value
  appliedFloorFilter.value = floorFilter.value
  appliedContractStatusFilter.value = contractStatusFilter.value
  appliedContractExpiryFilter.value = contractExpiryFilter.value
  appliedContractExpiryFrom.value = contractExpiryFrom.value
  appliedContractExpiryTo.value = contractExpiryTo.value
  appliedBookingRangeFilter.value = bookingRangeFilter.value
  appliedBookingFrom.value = bookingFrom.value
  appliedBookingTo.value = bookingTo.value
  appliedRenewalStatusFilter.value = renewalStatusFilter.value
  appliedReturnStatusFilter.value = returnStatusFilter.value
  appliedLegalStatusFilter.value = legalStatusFilter.value
  appliedMaintenanceStatusFilter.value = maintenanceStatusFilter.value
  appliedRequestTypeFilter.value = requestTypeFilter.value
  appliedCategoryFilter.value = draftCategoryFilter.value
  emit('update:categoryFilter', draftCategoryFilter.value || null)
  emit('update:onlyAvailable', availabilityFilter.value === 'available')
  showAdvancedFilters.value = false
  refreshSources()
  if (focusAppliedProjectFilterResults()) return
}

const toggleCategoryFilter = (catName: string) => {
  draftCategoryFilter.value = draftCategoryFilter.value === catName || !catName ? '' : catName
}

const handleCategorySelect = (e: Event) => {
  const target = e.target as HTMLSelectElement
  draftCategoryFilter.value = target.value || ''
}

// Watchers for reactive updates
watch(isSystemDark, () => {
  if (!map) return
  map.setStyle(getStyleConfig(resolveStyle(currentStyle.value)))
  map.once('styledata', () => {
    setupProjectLayers()
    refreshSources()
  })
})
watch(() => props.mapStyle, (newStyle) => {
  const target = newStyle === 'satellite' ? 'satellite' : newStyle === 'dark' ? 'dark' : 'standard'
  if (target !== currentStyle.value) {
    setMapStyle(target)
  }
})

watch(() => props.categoryFilter, (category) => {
  const value = category || ''
  draftCategoryFilter.value = value
  appliedCategoryFilter.value = value
  refreshSources()
})

watch(() => props.onlyAvailable, (onlyAvailable) => {
  // A quick filter set by the parent is an explicit action, so apply it now.
  if (onlyAvailable) {
    availabilityFilter.value = 'available'
    appliedAvailabilityFilter.value = 'available'
  } else if (availabilityFilter.value === 'available' || appliedAvailabilityFilter.value === 'available') {
    availabilityFilter.value = 'all'
    appliedAvailabilityFilter.value = 'all'
  }
  refreshSources()
})

watch(() => props.stallStatusFilter, (status) => {
  availabilityFilter.value = status
  appliedAvailabilityFilter.value = status
  refreshSources()
})

watch(() => props.statusFilterContext, (context) => {
  let changed = false
  if (context !== 'legal_case' && (legalStatusFilter.value !== 'all' || appliedLegalStatusFilter.value !== 'all')) {
    legalStatusFilter.value = 'all'
    appliedLegalStatusFilter.value = 'all'
    changed = true
  }
  if (context !== 'action_required') {
    if (maintenanceStatusFilter.value !== 'all' || appliedMaintenanceStatusFilter.value !== 'all') changed = true
    if (requestTypeFilter.value !== 'all' || appliedRequestTypeFilter.value !== 'all') changed = true
    maintenanceStatusFilter.value = 'all'
    appliedMaintenanceStatusFilter.value = 'all'
    requestTypeFilter.value = 'all'
    appliedRequestTypeFilter.value = 'all'
  }
  if (changed) refreshSources()
})

watch(() => props.actionRequiredStallCodes, () => {
  if (appliedAvailabilityFilter.value === 'action_required') {
    refreshSources()
  }
})

watch(() => props.activeFilterCodes, () => {
  refreshSources()
})

watch(() => props.operationalFilterMeta, () => {
  if (hasOperationalFilter(appliedOperationalFilters())) refreshSources()
})

watch([
  () => props.selectedStallCodes.join('\u001f'),
  () => props.selectedStallKeys.join('\u001f')
], () => {
  map?.getSource('project-stalls')?.setData(getStallsGeoJson())
  renderSelectedStallPins()
})

watch(() => props.marketFull, () => {
  refreshSources()
  renderSelectedStallPins()
})

watch(() => props.stallStatusOverrides, () => {
  clearProjectLayoutCache()
  scheduleStatusSourceRefresh()
})

// Parent explicitly bumps this key after a GIS edit/load. Avoid traversing the
// complete 13MB projects tree with Vue deep-watch on every nested mutation.
watch(() => props.managedDataVersion, () => {
  applyManagedProjects(props.managedProjects)
  applyManagedFacilities(props.managedFacilities)
})

// Parent may switch project externally (e.g. deep-link or breadcrumb reset); follow it if different
watch(() => props.selectedProjectId, (newId) => {
  if (!map) return
  if (!newId) {
    if (activeProjectId.value) exitProject()
    return
  }
  if (newId === activeProjectId.value) return
  if (props.managementMode && isPolygonEditMode.value && !activeProjectId.value) return
  const found = projects.value.find(p => p.id === newId)
  if (found) enterProject(found)
})

let marketSwitchController: AbortController | null = null
// Dynamic market switcher watcher: loads dedicated JSON file on demand and flies to market location
watch(() => props.marketId, async (newMarketId) => {
  marketSwitchController?.abort()
  const controller = new AbortController()
  marketSwitchController = controller
  activeProjectId.value = null
  hoveredProject.value = null
  hoveredStallKey.value = null
  currentFloorId.value = '1'
  try {
    await syncMarketJson(newMarketId, controller.signal)
  } catch (error: any) {
    if (controller.signal.aborted || error?.name === 'AbortError') return
    throw error
  } finally {
    if (marketSwitchController === controller) marketSwitchController = null
  }
  if (controller.signal.aborted) return
  const selectedProject = props.selectedProjectId
    ? projects.value.find(project => project.id === props.selectedProjectId)
    : null
  const preferredFloorId = selectedProject ? getProjectActualFloorIds(selectedProject)[0] : null
  selectedMarketFloorId.value = preferredFloorId || (marketFloors.value.length > 1 ? 'all' : (marketFloors.value[0]?.id || ''))
  marketFloorBeforeProject.value = marketFloors.value.length > 1 ? 'all' : selectedMarketFloorId.value
  if (map) {
    refreshSources()
    renderSelectedStallPins()
    map.flyTo({
      center: currentMarketCenter.value,
      zoom: currentMarketZoom.value,
      pitch: currentMarketPitch.value,
      bearing: currentMarketBearing.value,
      duration: 900
    })
  }
}, { immediate: true })

const resize = () => {
  map?.resize()
}

const finishActiveTools = () => {
  cancelDrawing()
  if (isStallPolygonEditMode.value) {
    emitCurrentEditingStallUpdate()
    isStallPolygonEditMode.value = false
    stallEditSnapshot.value = null
    stallEditHistory.value = []
    clearStallVertexHandles()
  }
  if (isPolygonEditMode.value) {
    isPolygonEditMode.value = false
    clearPolygonVertexHandles()
  }
  refreshSources()
}

// ---------------- หน้าข้อมูลตลาด: แผงสไลด์ ค้นหา / กรอง / ชั้นข้อมูล / แดชบอร์ด ----------------
const ovPanel = ref<'dashboard' | 'search' | 'filter' | 'layers' | null>(props.isDashboardOpen ? 'dashboard' : null)

watch(ovPanel, (val) => {
  emit('update:activePanel', val)
  emit('update:isDashboardOpen', val === 'dashboard')
  if (process.client) {
    window.setTimeout(() => { map?.resize() }, 280)
  }
}, { immediate: true })

const ovSearchInput = ref<HTMLInputElement | null>(null)
const ovTitle = computed(() => {
  if (ovPanel.value === 'dashboard') return props.dashboardTitle || 'แดชบอร์ด'
  if (ovPanel.value === 'search') return 'ค้นหา'
  if (ovPanel.value === 'filter') return props.filterContext === 'costs' ? 'ตัวกรองค่าใช้จ่าย' : 'ตัวกรอง'
  return 'ชั้นข้อมูล'
})
const ovSubtitle = computed(() => {
  if (ovPanel.value === 'dashboard') return props.dashboardSubtitle || 'สรุปข้อมูลภาพรวม'
  if (ovPanel.value === 'search') return 'ค้นหาการจอง คำร้อง หรือแผง'
  if (ovPanel.value === 'filter') return props.filterContext === 'costs' ? 'หาบิลที่ต้องติดตาม' : 'หาแผงและโซนที่ตรงเงื่อนไข'
  return 'เลือกแผนที่ฐานและดูความหมายของสี'
})
// กรองจากสีโครงการ (ชั้นข้อมูล) แยกจากตัวกรอง จึงไม่นับในตัวเลขบนปุ่มตัวกรอง
const ovFilterBadge = computed(() => (props.filterContext === 'costs'
  ? props.filterBadge
  : advancedFilterCount.value))
// Layer visibility is a display preference, not an active data filter.
// Only selected project-colour categories contribute to the badge count.
const mapLayerFilterCount = computed(() => legendCategories.value.length)
const basemapOptions: { value: 'standard' | 'satellite' | 'dark'; label: string }[] = [
  { value: 'standard', label: 'แผนที่ปกติ' },
  { value: 'satellite', label: 'ดาวเทียม' }
]
const legendCategoryStyle = (color: string): Record<string, string> => ({ '--ov-category-color': color })
const toggleOvPanel = (panel: 'dashboard' | 'search' | 'filter' | 'layers') => {
  ovPanel.value = ovPanel.value === panel ? null : panel
  // preventScroll: แผงกำลังสไลด์เข้าจากนอกจอ ถ้า focus แบบปกติเบราว์เซอร์จะเลื่อน container ไปหาช่องค้นหา ทำให้แผนที่เคลื่อนผิดปกติ
  if (ovPanel.value === 'search') nextTick(() => ovSearchInput.value?.focus({ preventScroll: true }))
}
const closeOvPanel = () => {
  ovPanel.value = null
}

watch(() => props.isDashboardOpen, (open) => {
  if (open && ovPanel.value !== 'dashboard') ovPanel.value = 'dashboard'
  else if (!open && ovPanel.value === 'dashboard') ovPanel.value = null
})
type SearchMatchType = 'booking' | 'application' | 'stall' | 'project'

const searchResults = computed(() => {
  const q = suggestionQuery.value.trim()
  if (!q) return []
  const normalizedQuery = q.toLowerCase()
  return overviewProjects.value.filter(project => matchesSearchWith(project, q)).slice(0, 30).map((project) => {
    const stallCodes = new Set<string>()
    const references = new Map<string, { type: 'booking' | 'application'; number: string }>()
    const matchedStalls: Array<{ floorId: string; stall: MapStall }> = []
    const matchedKeys = new Set<string>()
    getEffectiveFloorLayoutsForProject(project).forEach((layout, floorId) => {
      layout.stalls.forEach((stall) => {
        const code = String(stall.code || '').trim().toUpperCase()
        const meta = props.operationalFilterMeta[code]
        const mapNumber = getStallMapLabel(stall).toLowerCase()
        const matchedReferences = (meta?.searchReferences || [])
          .filter(reference => reference.number.toLowerCase().includes(normalizedQuery))
        const matchesVisibleNumber = mapNumber.includes(normalizedQuery)
        const matchesInternalCode = !/^\d+$/.test(normalizedQuery) && code.toLowerCase().includes(normalizedQuery)
        if (matchesVisibleNumber || matchesInternalCode || matchedReferences.length) {
          stallCodes.add(getStallMapLabel(stall) || String(stall.name || '').trim() || code)
          const key = `${floorId}:${getStallSelectionKey(stall)}`
          if (!matchedKeys.has(key)) {
            matchedKeys.add(key)
            matchedStalls.push({ floorId, stall })
          }
        }
        matchedReferences.forEach(reference => references.set(`${reference.type}:${reference.number}`, reference))
      })
    })
    const referenceLabels = Array.from(references.values()).map(reference =>
      `${reference.type === 'booking' ? 'เลขที่การจอง' : 'เลขที่คำร้อง'} ${reference.number}`
    )
    const referenceText = [referenceLabels.join(', '), stallCodes.size ? `แผง ${Array.from(stallCodes).join(', ')}` : '']
      .filter(Boolean).join(' · ')
    const matchTypes = new Set<SearchMatchType>()
    Array.from(references.values()).forEach(reference => matchTypes.add(reference.type))
    if (stallCodes.size) matchTypes.add('stall')
    if (!matchTypes.size) matchTypes.add('project')
    return { project, referenceText, matchedStalls, matchTypes: Array.from(matchTypes) }
  })
})
const searchResultStallCount = computed(() => searchResults.value.reduce((total, result) => total + result.matchedStalls.length, 0))
const searchMatchTypeLabel = (type: SearchMatchType) => ({
  booking: 'การจอง',
  application: 'คำร้อง',
  stall: 'เลขแผง',
  project: 'โครงการ'
}[type])
const getHighlightedSearchParts = (text: string) => {
  const query = suggestionQuery.value.trim()
  if (!query) return [{ text, match: false }]
  const parts: Array<{ text: string; match: boolean }> = []
  const normalizedText = text.toLocaleLowerCase('th-TH')
  const normalizedQuery = query.toLocaleLowerCase('th-TH')
  let cursor = 0
  let matchIndex = normalizedText.indexOf(normalizedQuery, cursor)
  while (matchIndex >= 0) {
    if (matchIndex > cursor) parts.push({ text: text.slice(cursor, matchIndex), match: false })
    parts.push({ text: text.slice(matchIndex, matchIndex + query.length), match: true })
    cursor = matchIndex + query.length
    matchIndex = normalizedText.indexOf(normalizedQuery, cursor)
  }
  if (cursor < text.length) parts.push({ text: text.slice(cursor), match: false })
  return parts.length ? parts : [{ text, match: false }]
}
const projectColor = (p: ProjectZone) => categories.value.find(c => c.name === p.category)?.color || '#94a3b8'
const moveSearchSuggestion = (direction: 1 | -1) => {
  const count = searchResults.value.length
  if (!count) return
  if (focusedSuggestionIndex.value < 0) focusedSuggestionIndex.value = direction > 0 ? 0 : count - 1
  else focusedSuggestionIndex.value = (focusedSuggestionIndex.value + direction + count) % count
  nextTick(() => {
    document.querySelector<HTMLElement>(`[data-suggestion-index="${focusedSuggestionIndex.value}"]`)
      ?.scrollIntoView({ block: 'nearest' })
  })
}
const selectActiveSearchSuggestion = () => {
  const result = searchResults.value[focusedSuggestionIndex.value >= 0 ? focusedSuggestionIndex.value : 0]
  if (result) selectSearchResult(result)
}
const selectSearchResult = (result: (typeof searchResults.value)[number]) => {
  selectedSuggestionId.value = String(result.project.id)
  focusedSuggestionIndex.value = searchResults.value.findIndex(item => String(item.project.id) === selectedSuggestionId.value)
  appliedSearchQuery.value = suggestionQuery.value.trim()
  refreshSources()
  if (result.matchedStalls.length) {
    const targetFloor = new Set(result.matchedStalls.map(match => match.floorId)).size === 1
      ? result.matchedStalls[0]?.floorId
      : undefined
    enterProject(result.project, targetFloor)
    nextTick(() => {
      const bounds = getStallsBounds(result.matchedStalls.map(match => match.stall))
      if (!bounds) return
      map?.fitBounds(bounds, {
        padding: getDrillMapPadding(),
        maxZoom: result.matchedStalls.length === 1 ? 21 : 19.5,
        pitch: DRILL_PITCH,
        bearing: map.getBearing(),
        duration: 800
      })
    })
    return
  }
  enterProject(result.project)
}

watch(searchResults, (results) => {
  if (!results.length) focusedSuggestionIndex.value = -1
  else if (focusedSuggestionIndex.value >= results.length) focusedSuggestionIndex.value = results.length - 1
})
const toggleLegendCategory = (name: string) => {
  legendCategories.value = legendCategories.value.includes(name)
    ? legendCategories.value.filter(n => n !== name)
    : [...legendCategories.value, name]
  refreshSources()
}
const clearLegendCategories = () => {
  legendCategories.value = []
  refreshSources()
}
const toggleProjectLabels = () => {
  showProjectLabels.value = !showProjectLabels.value
  renderZoneLabels()
}
const toggleStallLabels = () => {
  showStallNumberLabels.value = !showStallNumberLabels.value
  if (map?.getLayer('stalls-label')) {
    map.setLayoutProperty('stalls-label', 'visibility', showStallNumberLabels.value ? 'visible' : 'none')
  }
}
const toggleFacilitiesLayer = () => {
  showFacilities.value = !showFacilities.value
  pinnedFacilityPopup.value = false
  hoverPopup?.remove()
  refreshSources()
}
const applyFiltersFromPanel = () => {
  applyFilters()
  ovPanel.value = null
}
const onOvKeydown = (e: KeyboardEvent) => {
  if (e.key === 'Escape' && ovPanel.value) ovPanel.value = null
}
onMounted(() => { if (process.client) window.addEventListener('keydown', onOvKeydown) })
onBeforeUnmount(() => {
  if (process.client) window.removeEventListener('keydown', onOvKeydown)
  if (searchDebounceTimer) window.clearTimeout(searchDebounceTimer)
})

// กดตัวเลขสถานะบนแดชบอร์ดเพื่อกรองแผงทันที (ใช้ตัวกรองชุดเดียวกับหน้าตัวกรอง) ส่งค่าที่ไม่ระบุ = คงเดิม
const applyDashboardFilter = (patch: Partial<{ zone: string; contractStatus: ContractStatusFilter; renewalStatus: RenewalFilter; returnStatus: ReturnFilter; legalStatus: LegalFilter; maintenanceStatus: MaintenanceFilter }>) => {
  if (patch.zone !== undefined) { zoneFilter.value = patch.zone || 'all'; appliedZoneFilter.value = zoneFilter.value }
  if (patch.contractStatus !== undefined) { contractStatusFilter.value = patch.contractStatus; appliedContractStatusFilter.value = patch.contractStatus }
  if (patch.renewalStatus !== undefined) { renewalStatusFilter.value = patch.renewalStatus; appliedRenewalStatusFilter.value = patch.renewalStatus }
  if (patch.returnStatus !== undefined) { returnStatusFilter.value = patch.returnStatus; appliedReturnStatusFilter.value = patch.returnStatus }
  if (patch.legalStatus !== undefined) { legalStatusFilter.value = patch.legalStatus; appliedLegalStatusFilter.value = patch.legalStatus }
  if (patch.maintenanceStatus !== undefined) { maintenanceStatusFilter.value = patch.maintenanceStatus; appliedMaintenanceStatusFilter.value = patch.maintenanceStatus }
  refreshSources()
}

watch([appliedSearchQuery, appliedZoneFilter, appliedContractStatusFilter, appliedRenewalStatusFilter, appliedReturnStatusFilter, appliedLegalStatusFilter, appliedMaintenanceStatusFilter], () => {
  emit('operational-filter-change', {
    zone: appliedZoneFilter.value,
    search: appliedSearchQuery.value,
    contractStatus: appliedContractStatusFilter.value,
    renewalStatus: appliedRenewalStatusFilter.value,
    returnStatus: appliedReturnStatusFilter.value,
    legalStatus: appliedLegalStatusFilter.value,
    maintenanceStatus: appliedMaintenanceStatusFilter.value
  })
}, { immediate: true })

// กดแผงในแดชบอร์ด: ใช้กลไกค้นหาเดียวกับกล่องค้นหา (เข้าโครงการ ซูมไปที่แผง หรี่แผงอื่น) กดซ้ำเพื่อล้าง
const searchStallFromDashboard = async (code: string) => {
  const query = String(code || '').trim()
  if (!query) return
  searchQuery.value = query
  suggestionQuery.value = query
  await nextTick()
  const result = searchResults.value[0]
  if (result) selectSearchResult(result)
}

defineExpose({ enterProject, exitProject, resize, finishActiveTools, focusSelectedStalls, focusBookedStalls, applyDashboardFilter, searchStallFromDashboard, clearAppliedSearch })

const canTeleportToolbar = ref(false)
let isUnmounted = false
let mapLoadController: AbortController | null = null

const handleWindowResize = () => {
  map?.resize()
}

const initializeMap = async () => {
  if (!process.client || map || isUnmounted) return
  mapLoadController?.abort()
  const controller = new AbortController()
  mapLoadController = controller
  canTeleportToolbar.value = Boolean(document.getElementById('step1-map-toolbar'))
  window.removeEventListener('resize', handleWindowResize)
  window.addEventListener('resize', handleWindowResize)
  try {
    // Start both jobs together and render the basemap as soon as MapLibre is
    // available. Large GIS JSON continues hydrating the sources in parallel.
    isMapBaseReady = false
    isGisHydrating = true
    const marketDataPromise = syncMarketJson(props.marketId, controller.signal).then(
      () => null,
      error => error
    ).finally(() => {
      if (controller.signal.aborted) return
      isGisHydrating = false
      revealMapIfReady()
    })
    const ml = await import('maplibre-gl')
    if (isUnmounted || controller.signal.aborted) return
    maplibregl = ml.default || ml
    initMap()
    const marketDataError = await marketDataPromise
    if (marketDataError) throw marketDataError
  } catch (error: any) {
    if (controller.signal.aborted || error?.name === 'AbortError') return
    // Do not leave a permanent loading surface if the map library or data
    // cannot be prepared; the existing page controls remain usable.
    console.error('[map] unable to initialize:', error)
    isMapLoading.value = false
  } finally {
    if (mapLoadController === controller) mapLoadController = null
  }
}

onMounted(() => {
  isUnmounted = false
  void initializeMap()
})

// วัดตำแหน่งขอบล่างของชุดปุ่มซูมจริง แล้วส่งเป็นตัวแปร --ov-ctrl-bottom ให้แถบเลือกชั้นวางขอบล่างเสมอกันเป๊ะ
let ctrlBaselineObserver: ResizeObserver | null = null
const watchCtrlBaseline = (navGroup: HTMLElement) => {
  const wrapper = navGroup.closest('.maplibre-overview-wrapper') as HTMLElement | null
  if (!wrapper || typeof ResizeObserver === 'undefined') return
  const sync = () => {
    const w = wrapper.getBoundingClientRect()
    const g = navGroup.getBoundingClientRect()
    if (!g.height) return
    wrapper.style.setProperty('--ov-ctrl-bottom', `${Math.max(0, Math.round(w.bottom - g.bottom))}px`)
  }
  ctrlBaselineObserver?.disconnect()
  ctrlBaselineObserver = new ResizeObserver(sync)
  ctrlBaselineObserver.observe(wrapper)
  const container = navGroup.parentElement?.parentElement
  if (container) ctrlBaselineObserver.observe(container)
  requestAnimationFrame(sync)
}
const cleanupMapResources = () => {
  ctrlBaselineObserver?.disconnect()
  ctrlBaselineObserver = null
  isUnmounted = true
  mapLoadController?.abort()
  mapLoadController = null
  marketSwitchController?.abort()
  marketSwitchController = null
  isMapLoading.value = true
  canTeleportToolbar.value = false
  if (statusRefreshFrame !== null) {
    window.cancelAnimationFrame(statusRefreshFrame)
    statusRefreshFrame = null
  }
  statusRefreshQueuedForMove = false
  if (process.client) {
    window.removeEventListener('resize', handleWindowResize)
  }
  try {
    hoverPopup?.remove()
  } catch {}
  selectedStallMarkers.forEach((m) => {
    try { m.remove() } catch {}
  })
  selectedStallMarkers = []
  clearPolygonVertexHandles()
  clearStallVertexHandles()
  clearFacilityMarkers()
  if (map) {
    const mapInstance = map
    map = null
    // Removing a large WebGL map synchronously blocks route navigation. Detach
    // it from Vue now, then release GPU resources once the next tab has painted.
    try { mapInstance.stop() } catch {}
    const releaseMap = () => {
      try { mapInstance.remove() } catch {}
    }
    if (process.client && typeof (window as any).requestIdleCallback === 'function') {
      ;(window as any).requestIdleCallback(releaseMap, { timeout: 500 })
    } else if (process.client) {
      window.setTimeout(releaseMap, 0)
    } else {
      releaseMap()
    }
  }
}

onBeforeUnmount(cleanupMapResources)

// A kept-alive GIS page keeps its WebGL context and parsed sources. Pause animation
// and global listeners while another page is active, then resize the existing map
// when the user returns instead of destroying and rebuilding MapLibre.
onDeactivated(() => {
  mapLoadController?.abort()
  mapLoadController = null
  marketSwitchController?.abort()
  marketSwitchController = null
  if (process.client) window.removeEventListener('resize', handleWindowResize)
  if (statusRefreshFrame !== null) {
    window.cancelAnimationFrame(statusRefreshFrame)
    statusRefreshFrame = null
  }
  statusRefreshQueuedForMove = false
  try { map?.stop() } catch {}
  try { hoverPopup?.remove() } catch {}
})

onActivated(() => {
  if (!process.client || isUnmounted) return
  window.removeEventListener('resize', handleWindowResize)
  window.addEventListener('resize', handleWindowResize)
  if (!map) {
    void initializeMap()
    return
  }
  nextTick(() => window.requestAnimationFrame(() => map?.resize()))
})
</script>

<style scoped>
.maplibre-overview-wrapper.is-drawing-active :deep(.maplibre-zone-number-badge),.maplibre-overview-wrapper.is-drawing-active :deep(.facility-map-marker){pointer-events:none!important}
.zone-management-toolbar{position:absolute;z-index:26;left:50%;bottom:18px;transform:translateX(-50%);display:flex;max-width:calc(100% - 32px);overflow-x:auto;background:var(--bg-rgba-255-255-255-0p97);border:1px solid var(--bd-d7e2ef);border-radius:14px;box-shadow:0 8px 20px rgba(15,23,42,.13);backdrop-filter:blur(8px)}.zone-management-toolbar :deep(.app-btn){height:68px;min-height:68px;max-height:68px;min-width:104px;padding:7px 14px;flex:0 0 auto;flex-direction:column;gap:5px;border:0!important;border-right:1px solid var(--bd-e8eef6)!important;border-radius:0!important;background:transparent!important;color:var(--tx-334155)!important;font-size:.8125rem;font-weight:600;line-height:1.15;white-space:nowrap;box-shadow:none!important}.zone-management-toolbar :deep(.app-btn:first-child){border-radius:13px 0 0 13px!important}.zone-management-toolbar :deep(.app-btn:last-child){border-right:0!important;border-radius:0 13px 13px 0!important}.zone-management-toolbar :deep(.app-btn:hover:not(.is-disabled)){background:var(--bg-f8fbff)!important;color:var(--tx-2563eb)!important}.zone-management-toolbar :deep(.app-btn.app-btn-primary:not(.is-disabled)){background:var(--bg-eff6ff)!important;color:var(--tx-2563eb)!important;font-weight:700!important}.zone-management-toolbar :deep(.app-btn.app-btn-primary:hover:not(.is-disabled)){background:var(--bg-dbeafe)!important;color:var(--tx-1d4ed8)!important}.zone-management-toolbar :deep(.app-btn.app-btn-danger-outline:not(.is-disabled)){color:var(--tx-dc2626)!important}.zone-management-toolbar :deep(.app-btn.app-btn-danger-outline:hover:not(.is-disabled)){background:var(--bg-fef2f2)!important;color:var(--tx-b91c1c)!important}.zone-management-toolbar :deep(.app-btn.is-disabled){background:var(--bg-ffffff)!important;color:var(--tx-cbd5e1)!important;opacity:1!important}.zone-management-toolbar :deep(.app-btn svg){width:20px;height:20px}.zone-management-toolbar>button:not(.app-btn){min-width:88px;padding:11px 13px;border:0;border-radius:8px;background:transparent;color:var(--tx-334155);display:grid;justify-items:center;gap:4px;font:600 .8125rem/1 inherit;cursor:pointer;white-space:nowrap}.zone-management-toolbar>button:not(.app-btn):hover:not(:disabled),.zone-management-toolbar>button:not(.app-btn).active{background:var(--bg-eff6ff);color:var(--tx-2563eb)}.zone-management-toolbar>button.primary-save:not(:disabled){background:#2563eb;color:#fff;padding-inline:18px}.zone-management-toolbar>button.primary-save:hover:not(:disabled){background:#1d4ed8;color:#fff}.zone-management-toolbar>button:disabled{opacity:.38;cursor:not-allowed}.zone-management-toolbar>button.danger:hover:not(:disabled){background:var(--bg-fef2f2);color:var(--tx-dc2626)}.polygon-edit-hint{position:absolute;z-index:26;left:50%;bottom:96px;transform:translateX(-50%);display:flex;align-items:center;gap:8px;padding:8px 16px;border-radius:999px;background:var(--bg-rgba-255-255-255-0p97);border:1px solid var(--bd-d7e2ef);color:var(--tx-1e293b);font-size:.8125rem;font-weight:600;box-shadow:0 8px 20px rgba(15,23,42,.12);backdrop-filter:blur(8px);white-space:nowrap}.polygon-edit-hint svg{width:16px;height:16px;color:var(--tx-2563eb);flex-shrink:0}@media(max-width:640px){.zone-management-toolbar{bottom:10px;left:16px;right:16px;transform:none}.zone-management-toolbar :deep(.app-btn){min-width:88px;height:62px;min-height:62px;max-height:62px;padding:6px 8px;font-size:.75rem}.zone-management-toolbar>button:not(.app-btn){min-width:76px;padding:9px 8px;font-size:.75rem}.polygon-edit-hint{bottom:82px;font-size:.75rem;white-space:normal;text-align:center}}
.projects-master-container {
  display: flex;
  flex-direction: column;
  gap: 16px;
  background: transparent;
  border-radius: 0;
  padding: 0;
}

/* Top bar with map mode switcher on the left and search on the right */
.master-map-top-bar {
  display: flex;
  justify-content: flex-end;
  align-items: center;
  gap: 12px;
  margin-bottom: 0;
  width: 100%;
  flex-wrap: wrap;
}

.master-map-tools-left {
  display: flex;
  align-items: center;
  gap: 10px;
  flex-wrap: wrap;
}

.master-map-tools-right {
  display: flex;
  align-items: center;
  justify-content: flex-end;
  gap: 9px;
  margin-left: auto;
  flex-wrap: wrap;
}


/* Filter Pill Button: เฉพาะแผงว่าง */
.map-filter-pill-btn {
  display: inline-flex;
  align-items: center;
  gap: 7px;
  background: var(--bg-ffffff);
  border: 1px solid var(--bd-cbd5e1);
  border-radius: 9999px;
  padding: 6px 14px;
  font-family: 'Bai Jamjuree', sans-serif;
  font-size: 0.84rem;
  font-weight: 600;
  color: var(--tx-475569);
  cursor: pointer;
  transition: all 0.2s cubic-bezier(0.16, 1, 0.3, 1);
  box-shadow: 0 1px 2px rgba(0, 0, 0, 0.03);
}

.map-filter-pill-btn:hover {
  border-color: var(--bd-94a3b8);
  color: var(--tx-1e293b);
  transform: translateY(-1px);
}

.map-filter-pill-btn.active {
  background: var(--bg-ecfdf5);
  border-color: #10b981;
  color: var(--tx-065f46);
  font-weight: 700;
  box-shadow: 0 2px 6px rgba(16, 185, 129, 0.15);
}

.filter-status-dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  background: #94a3b8;
  transition: all 0.2s ease;
}

.filter-status-dot.is-active {
  background: #10b981;
  box-shadow: 0 0 0 2px rgba(16, 185, 129, 0.25);
}

.filter-count-badge {
  font-size: 0.76rem;
  font-weight: 700;
  color: inherit;
  opacity: 0.85;
}

.map-advanced-filter-wrap { position: relative; }
.map-advanced-filter-btn { display: inline-flex; align-items: center; gap: 6px; height: 35px; padding: 6px 12px; border: 1px solid var(--bd-cbd5e1); border-radius: 8px; background: var(--bg-ffffff); color: var(--tx-475569); font-family: inherit; font-size: .82rem; font-weight: 700; cursor: pointer; }
.map-advanced-filter-btn:hover, .map-advanced-filter-btn.active { border-color: var(--bd-93c5fd); background: var(--bg-eff6ff); color: var(--tx-1d4ed8); }
.map-advanced-count { display: grid; width: 17px; height: 17px; place-items: center; border-radius: 50%; background: #2563eb; color: #fff; font-size: .65rem; }
.map-advanced-filter-panel {
  position: absolute;
  top: calc(100% + 8px);
  right: 0;
  left: auto;
  z-index: 45;
  width: min(320px, calc(100vw - 64px));
  padding: 16px;
  border: 1px solid var(--bd-dbe5f2);
  border-radius: 14px;
  background: var(--bg-ffffff);
  box-shadow: 0 16px 36px rgba(15, 23, 42, 0.16), 0 4px 12px rgba(15, 23, 42, 0.05);
}
.map-filter-panel-head {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 6px;
  color: var(--tx-0f172a);
  font-size: 0.88rem;
  font-weight: 700;
}
.map-filter-panel-head button {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 26px;
  height: 26px;
  border: 0;
  border-radius: 6px;
  background: transparent;
  color: var(--tx-64748b);
  cursor: pointer;
  transition: all 0.15s ease;
}
.map-filter-panel-head button:hover {
  background: var(--bg-f1f5f9);
  color: var(--tx-0f172a);
}
.filter-field-group {
  margin-top: 10px;
}
.filter-field-label {
  display: block;
  margin-bottom: 5px;
  color: var(--tx-475569);
  font-size: 0.78rem;
  font-weight: 600;
}
.map-advanced-filter-panel :deep(.app-select-trigger) {
  min-height: 38px;
  height: 38px;
  padding: 0 12px;
  font-size: 0.84rem;
  border-radius: 8px;
}
.map-advanced-filter-panel :deep(.app-select-menu) {
  z-index: 100;
  max-height: 220px;
}
.map-advanced-filter-panel :deep(.app-select-option) {
  padding: 8px 10px;
  font-size: 0.82rem;
}
.map-filter-panel-footer {
  display: flex;
  flex-direction: column;
  gap: 12px;
  margin-top: 14px;
  padding-top: 12px;
  border-top: 1px solid var(--bd-f1f5f9);
}

.mfp-result-summary {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 6px;
  font-size: 0.78rem;
  color: var(--tx-64748b);
  background: var(--bg-f8fafc);
  padding: 6px 10px;
  border-radius: 8px;
  border: 1px solid var(--bd-eef2f6);
}

.mfp-result-label {
  color: var(--tx-94a3b8);
  font-size: 0.75rem;
}

.mfp-result-val strong {
  color: var(--tx-2563eb);
  font-weight: 700;
}

.mfp-result-dot {
  color: var(--tx-cbd5e1);
}

.mfp-result-stalls strong {
  color: var(--tx-0f172a);
  font-weight: 600;
}

.mfp-action-grid {
  display: grid;
  grid-template-columns: 1fr 1.6fr;
  gap: 8px;
  width: 100%;
}

.btn-filter-reset {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: 6px;
  padding: 8px 12px;
  border-radius: 8px;
  border: 1.5px solid var(--bd-cbd5e1);
  background: var(--bg-ffffff);
  color: var(--tx-334155);
  font-family: inherit;
  font-size: 0.82rem;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.16s ease;
}

.btn-filter-reset:hover {
  background: var(--bg-f8fafc);
  border-color: var(--bd-94a3b8);
  color: var(--tx-0f172a);
}

.btn-filter-apply {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: 6px;
  padding: 8px 14px;
  border-radius: 8px;
  border: none;
  background: #2563eb;
  color: #ffffff;
  font-family: inherit;
  font-size: 0.82rem;
  font-weight: 700;
  cursor: pointer;
  box-shadow: 0 2px 6px rgba(37, 99, 235, 0.25);
  transition: all 0.16s ease;
}

.btn-filter-apply:hover {
  background: #1d4ed8;
  box-shadow: 0 4px 10px rgba(37, 99, 235, 0.35);
  transform: translateY(-1px);
}

/* Category Select Dropdown */
.map-category-select-wrapper {
  position: relative;
  display: inline-flex;
  align-items: center;
}

.map-category-select {
  appearance: none;
  -webkit-appearance: none;
  background: var(--bg-ffffff);
  border: 1px solid var(--bd-cbd5e1);
  border-radius: 9999px;
  padding: 6px 32px 6px 14px;
  font-family: 'Bai Jamjuree', sans-serif;
  font-size: 0.84rem;
  font-weight: 600;
  color: var(--tx-334155);
  cursor: pointer;
  transition: all 0.2s ease;
  box-shadow: 0 1px 2px rgba(0, 0, 0, 0.03);
  outline: none;
  height: 35px;
}

.map-category-select:hover {
  border-color: var(--bd-94a3b8);
  color: var(--tx-0f172a);
}

.map-category-select:focus {
  border-color: #2563eb;
  box-shadow: 0 0 0 3px rgba(37, 99, 235, 0.12);
}

.map-filter-select {
  height: 35px;
  max-width: 175px;
  padding: 6px 9px;
  border: 1px solid var(--bd-cbd5e1);
  border-radius: 8px;
  outline: none;
  background: var(--bg-ffffff);
  color: var(--tx-475569);
  font-family: 'Bai Jamjuree', sans-serif;
  font-size: .78rem;
  font-weight: 600;
  cursor: pointer;
}
.map-filter-select:hover { border-color: var(--bd-94a3b8); color: var(--tx-1e293b); }
.map-filter-select:focus { border-color: #2563eb; box-shadow: 0 0 0 3px rgba(37, 99, 235, .1); }
.map-filter-result { color: var(--tx-64748b); font-size: .76rem; white-space: nowrap; }
.map-filter-result strong { color: var(--tx-1e40af); font-weight: 800; }
.map-clear-filters-btn { border: 0; padding: 4px 2px; background: transparent; color: var(--tx-2563eb); font-family: inherit; font-size: .76rem; font-weight: 700; cursor: pointer; white-space: nowrap; }
.map-clear-filters-btn:hover { color: var(--tx-1d4ed8); text-decoration: underline; }

.category-select-icon {
  position: absolute;
  right: 12px;
  pointer-events: none;
  color: var(--tx-64748b);
  display: flex;
  align-items: center;
  justify-content: center;
}

/* 2-Mode Segmented Control: แผนที่ vs ดาวเทียม */
.map-mode-segmented {
  display: inline-flex;
  align-items: center;
  background: var(--bg-f1f5f9);
  border: 1px solid var(--bd-e2e8f0);
  border-radius: 9999px;
  padding: 3px;
  gap: 2px;
  box-shadow: 0 1px 2px rgba(0, 0, 0, 0.03);
}

.map-mode-tab {
  border: none;
  background: transparent;
  color: var(--tx-64748b);
  font-family: 'Bai Jamjuree', sans-serif;
  font-size: 0.85rem;
  font-weight: 600;
  padding: 6px 18px;
  border-radius: 9999px;
  cursor: pointer;
  transition: all 0.2s cubic-bezier(0.16, 1, 0.3, 1);
  white-space: nowrap;
}

.map-mode-tab:hover {
  color: var(--tx-1e293b);
}

.map-mode-tab.active {
  background: var(--bg-ffffff);
  color: var(--tx-2563eb);
  font-weight: 700;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.08), 0 1px 2px rgba(0, 0, 0, 0.04);
}

.search-input-group {
  position: relative;
  display: inline-flex;
  align-items: center;
  width: 320px;
  max-width: 100%;
}

@media (max-width: 680px) {
  .master-map-tools-right { width: 100%; margin-left: 0; }
  .search-input-group { flex: 1; width: auto; }
}

.search-icon {
  position: absolute;
  left: 12px;
  top: 50%;
  transform: translateY(-50%);
  color: var(--tx-94a3b8);
  display: flex;
  align-items: center;
  justify-content: center;
  pointer-events: none;
}

.map-search-input {
  width: 100%;
  height: 38px;
  padding: 8px 32px 8px 36px;
  border-radius: 8px;
  border: 1px solid var(--bd-e2e8f0);
  background: var(--bg-ffffff);
  font-size: 0.85rem;
  font-family: inherit;
  color: var(--tx-1e293b);
  outline: none;
  transition: border-color 0.15s ease, box-shadow 0.15s ease;
}

.map-search-input:focus {
  border-color: #2563eb;
  box-shadow: 0 0 0 3px rgba(37, 99, 235, 0.08);
}

.map-search-input::placeholder {
  color: var(--tx-94a3b8);
}

.clear-search-btn {
  position: absolute;
  right: 10px;
  top: 50%;
  transform: translateY(-50%);
  background: none;
  border: none;
  color: var(--tx-94a3b8);
  cursor: pointer;
  padding: 4px;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: color 0.15s ease;
}

.clear-search-btn:hover {
  color: var(--tx-475569);
}


/* MapLibre Canvas Container */
.maplibre-overview-wrapper {
  position: relative;
  width: 100%;
  height: 620px;
  border-radius: 14px;
  overflow: hidden;
  border: 1px solid var(--bd-e2e8f0);
  box-shadow: 0 4px 16px rgba(15, 23, 42, 0.05);
}

.maplibre-overview-canvas {
  width: 100%;
  height: 100%;
}

/* Preserve the final map's spatial hierarchy while MapLibre fetches its style
   and tiles. The overlay is intentionally non-interactive so it never traps
   focus or pointer events once the map is ready. */
.map-loading-skeleton {
  position: absolute;
  inset: 0;
  z-index: 30;
  overflow: hidden;
  pointer-events: none;
  background:
    linear-gradient(120deg, var(--bg-rgba-226-232-240-0p5) 12%, var(--bg-rgba-248-250-252-0p96) 30%, var(--bg-rgba-226-232-240-0p5) 48%) 0 0 / 260% 100%,
    var(--bg-f8fafc);
  animation: mapSkeletonSweep 1.7s ease-in-out infinite;
}

.map-loading-skeleton::before,
.map-loading-skeleton::after {
  content: '';
  position: absolute;
  inset: auto;
  border: 1px solid var(--bd-rgba-203-213-225-0p72);
  background: var(--bg-rgba-226-232-240-0p36);
}

.map-loading-skeleton::before {
  width: 54%;
  height: 18%;
  right: -7%;
  top: 35%;
  transform: rotate(-12deg);
}

.map-loading-skeleton::after {
  width: 44%;
  height: 23%;
  left: -9%;
  bottom: 9%;
  transform: rotate(17deg);
}

.map-skeleton-control-group {
  position: absolute;
  z-index: 1;
  top: 18px;
  right: 18px;
  display: flex;
  align-items: center;
  gap: 8px;
}

.map-skeleton-line,
.map-skeleton-tool,
.map-skeleton-zone {
  display: block;
  border: 1px solid var(--bd-rgba-203-213-225-0p78);
  background: var(--bg-rgba-226-232-240-0p75);
}

.is-map-tab { width: 68px; height: 34px; border-radius: 999px; }
.is-map-tab.short { width: 56px; }
.is-center-control { width: 112px; height: 34px; border-radius: 999px; }

.map-skeleton-zone {
  position: absolute;
  z-index: 1;
  border-radius: 8px;
  background: var(--bg-rgba-191-219-254-0p48);
  box-shadow: inset 0 0 0 1px rgba(147, 197, 253, .22);
}

.map-skeleton-zone.zone-one { width: 18%; height: 25%; top: 27%; left: 21%; transform: rotate(-8deg); }
.map-skeleton-zone.zone-two { width: 15%; height: 36%; top: 35%; left: 45%; transform: rotate(9deg); }
.map-skeleton-zone.zone-three { width: 22%; height: 20%; right: 14%; bottom: 20%; transform: rotate(-13deg); }

.map-skeleton-legend {
  position: absolute;
  z-index: 1;
  left: 18px;
  bottom: 18px;
  width: 184px;
  padding: 13px;
  border: 1px solid var(--bd-rgba-203-213-225-0p76);
  border-radius: 12px;
  background: var(--bg-rgba-255-255-255-0p78);
}

.is-legend-title { width: 74px; height: 12px; border-radius: 999px; margin-bottom: 10px; }
.is-legend-row { width: 100%; height: 10px; border-radius: 999px; margin-top: 7px; }
.is-legend-row.shorter { width: 68%; }

.map-skeleton-tools {
  position: absolute;
  z-index: 1;
  left: 50%;
  bottom: 18px;
  display: flex;
  gap: 8px;
  transform: translateX(-50%);
}

.map-skeleton-tool { width: 82px; height: 48px; border-radius: 10px; background: var(--bg-rgba-255-255-255-0p82); }

@keyframes mapSkeletonSweep {
  0%, 100% { background-position: 100% 0; }
  50% { background-position: 0 0; }
}

@media (prefers-reduced-motion: reduce) {
  .map-loading-skeleton { animation: none; }
}

.map-project-guidance {
  position: absolute;
  top: 18px;
  left: 18px;
  z-index: 10;
  display: flex;
  align-items: center;
  gap: 10px;
  max-width: min(360px, calc(100% - 160px));
  padding: 11px 14px;
  border: 1px solid var(--bd-bfdbfe);
  border-radius: 12px;
  background: var(--bg-rgba-255-255-255-0p96);
  box-shadow: 0 8px 22px rgba(15, 23, 42, 0.12);
  color: var(--tx-1e3a8a);
  pointer-events: none;
}

.map-project-guidance-icon {
  display: grid;
  width: 34px;
  height: 34px;
  place-items: center;
  flex: 0 0 auto;
  border-radius: 9px;
  background: #2563eb;
  color: #ffffff;
  box-shadow: 0 4px 10px rgba(37, 99, 235, 0.25);
}

.map-project-guidance > span:nth-child(2) {
  display: grid;
  gap: 1px;
}

.map-project-guidance strong {
  font-size: 0.9rem;
  font-weight: 800;
  color: var(--tx-0f172a);
}

.map-project-guidance small {
  font-size: 0.75rem;
  color: var(--tx-64748b);
  line-height: 1.35;
}

.map-project-guidance-arrow {
  flex: 0 0 auto;
  margin-left: auto;
  color: var(--tx-2563eb);
}

@media (max-width: 680px) {
  .map-project-guidance {
    top: 12px;
    left: 12px;
    max-width: calc(100% - 24px);
  }
}

/* Floating Map Layer Legend (inside map - Bottom Left) */
.map-floating-legend-container {
  display: none !important;
  position: absolute;
  bottom: 16px;
  left: 16px;
  z-index: 20;
  display: flex;
  flex-direction: column;
  align-items: flex-start;
}

/* Floating Capsule Button - สีน้ำเงิน Royal Blue ตัดขอบขาว โดดเด่นสะดุดตาบนแผนที่ */
.map-floating-legend-btn {
  display: inline-flex;
  align-items: center;
  gap: 9px;
  background: #2563eb;
  color: #ffffff;
  padding: 7px 16px 7px 12px;
  border-radius: 9999px;
  border: 2px solid var(--bd-ffffff);
  box-shadow: 0 4px 18px rgba(37, 99, 235, 0.45), 0 2px 6px rgba(15, 23, 42, 0.2);
  font-family: inherit;
  font-size: 0.84rem;
  font-weight: 700;
  cursor: pointer;
  transition: all 0.2s cubic-bezier(0.16, 1, 0.3, 1);
  user-select: none;
}

.map-floating-legend-btn:hover {
  transform: translateY(-2px) scale(1.02);
  background: #1d4ed8;
  box-shadow: 0 8px 24px rgba(37, 99, 235, 0.55), 0 4px 10px rgba(15, 23, 42, 0.25);
}

.map-floating-legend-btn.active {
  background: #1e40af;
  border-color: var(--bd-93c5fd);
  box-shadow: 0 6px 20px rgba(30, 64, 175, 0.6);
}

.legend-btn-icon-wrap {
  width: 24px;
  height: 24px;
  border-radius: 50%;
  background: var(--bg-rgba-255-255-255-0p22);
  color: #ffffff;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  transition: all 0.15s ease;
}

.map-floating-legend-btn:hover .legend-btn-icon-wrap {
  background: var(--bg-rgba-255-255-255-0p32);
}

.legend-btn-label {
  letter-spacing: -0.01em;
  color: #ffffff;
}

.legend-count-pill {
  background: var(--bg-ffffff);
  color: var(--tx-2563eb);
  font-size: 0.76rem;
  font-weight: 800;
  padding: 1.5px 7.5px;
  border-radius: 9999px;
  min-width: 22px;
  text-align: center;
  line-height: 1.2;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.15);
}

.legend-chevron {
  transition: transform 0.2s cubic-bezier(0.16, 1, 0.3, 1);
  color: #ffffff;
}

.legend-chevron.is-open {
  transform: rotate(180deg);
}

/* Floating Popover Card */
.map-floating-legend-card {
  position: absolute;
  bottom: calc(100% + 12px);
  left: 0;
  width: 440px;
  max-width: calc(100vw - 32px);
  max-height: 380px;
  background: var(--bg-ffffff);
  border: 1.5px solid var(--bd-cbd5e1);
  border-radius: 16px;
  box-shadow: 0 20px 42px rgba(15, 23, 42, 0.22);
  display: flex;
  flex-direction: column;
  overflow: hidden;
  z-index: 25;
}

.mflc-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px 16px;
  border-bottom: 1px solid var(--bd-f1f5f9);
  background: var(--bg-ffffff);
}

.mflc-title-group {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  color: var(--tx-1e293b);
  font-size: 0.88rem;
}

.mflc-title-group svg {
  color: var(--tx-2563eb);
}

.mflc-sub {
  font-size: 0.75rem;
  color: var(--tx-64748b);
  background: var(--bg-f1f5f9);
  padding: 2px 7px;
  border-radius: 6px;
  font-weight: 500;
}

.mflc-close-btn {
  background: transparent;
  border: none;
  width: 28px;
  height: 28px;
  border-radius: 8px;
  display: grid;
  place-items: center;
  color: var(--tx-94a3b8);
  cursor: pointer;
  transition: all 0.15s ease;
}

.mflc-close-btn:hover {
  background: var(--bg-fee2e2);
  color: var(--tx-dc2626);
}

.mflc-body {
  padding: 12px 16px;
  overflow-y: auto;
  max-height: 300px;
}

.mflc-grid {
  list-style: none;
  margin: 0;
  padding: 0;
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 12px 14px;
}

.mflc-item {
  display: flex;
  align-items: flex-start;
  gap: 8px;
  font-size: 0.78rem;
  padding: 4px;
  border-radius: 8px;
  transition: background 0.15s ease;
}

.mflc-item:hover {
  background: var(--bg-f8fafc);
}

.mflc-swatch {
  flex-shrink: 0;
  width: 14px;
  height: 14px;
  margin-top: 2px;
  border-radius: 4px;
  border: 1.5px solid;
  opacity: 0.9;
}

.mflc-text {
  display: flex;
  flex-direction: column;
  min-width: 0;
}

.mflc-name {
  font-size: 0.78rem;
  font-weight: 600;
  color: var(--tx-1e293b);
  line-height: 1.35;
}

.mflc-projects {
  font-size: 0.75rem;
  color: var(--tx-64748b);
  margin-top: 2px;
  line-height: 1.3;
}

/* Transitions */
.legend-popover-fade-enter-active,
.legend-popover-fade-leave-active {
  transition: all 0.2s cubic-bezier(0.16, 1, 0.3, 1);
}

.legend-popover-fade-enter-from,
.legend-popover-fade-leave-to {
  opacity: 0;
  transform: translateY(8px) scale(0.97);
}

/* Modern Category Filter Panel */
.projects-category-filter-panel {
  background: var(--bg-ffffff);
  border: 1px solid var(--bd-e2e8f0);
  border-radius: 14px;
  padding: 16px 18px;
  box-shadow: 0 1px 3px rgba(15, 23, 42, 0.03);
  margin-top: 14px;
}

.filter-panel-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  flex-wrap: wrap;
  gap: 10px;
  margin-bottom: 14px;
  padding-bottom: 12px;
  border-bottom: 1px solid var(--bd-f1f5f9);
}

.filter-header-left {
  display: inline-flex;
  align-items: center;
  gap: 10px;
  flex-wrap: wrap;
}

.filter-title {
  font-size: 0.88rem;
  font-weight: 700;
  color: var(--tx-0f172a);
  display: inline-flex;
  align-items: center;
  gap: 7px;
}

.filter-title svg {
  color: var(--tx-2563eb);
}

.filter-subtitle-badge {
  font-size: 0.75rem;
  font-weight: 600;
  color: var(--tx-64748b);
  background: var(--bg-f1f5f9);
  padding: 3px 9px;
  border-radius: 999px;
  border: 1px solid var(--bd-e2e8f0);
  transition: all 0.2s ease;
}

.filter-subtitle-badge.is-filtering {
  background: var(--bg-eff6ff);
  border-color: var(--bd-bfdbfe);
  color: var(--tx-1e40af);
}

.btn-reset-filter {
  display: inline-flex;
  align-items: center;
  gap: 5px;
  padding: 4px 10px;
  border-radius: 8px;
  border: 1px solid var(--bd-e2e8f0);
  background: var(--bg-ffffff);
  color: var(--tx-dc2626);
  font-size: 0.75rem;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.15s ease;
}

.btn-reset-filter:hover {
  background: var(--bg-fef2f2);
  border-color: var(--bd-fca5a5);
}

/* Chips Flow Layout */
.category-chips-flow {
  display: flex;
  flex-wrap: wrap;
  gap: 8px 10px;
}

.category-filter-chip {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  padding: 6px 13px;
  border-radius: 999px;
  background: var(--bg-ffffff);
  border: 1px solid var(--bd-e2e8f0);
  color: var(--tx-334155);
  font-size: 0.8rem;
  font-weight: 500;
  box-shadow: 0 1px 2px rgba(15, 23, 42, 0.03);
  cursor: pointer;
  transition: all 0.2s ease;
  user-select: none;
  font-family: inherit;
}

.category-filter-chip:hover {
  border-color: var(--bd-cbd5e1);
  background: var(--bg-f8fafc);
  color: var(--tx-0f172a);
  transform: translateY(-1px);
  box-shadow: 0 3px 6px rgba(15, 23, 42, 0.05);
}

.category-filter-chip.active-chip {
  background: var(--bg-eff6ff);
  border-color: #2563eb;
  color: var(--tx-1e40af);
  font-weight: 700;
  box-shadow: 0 2px 8px rgba(37, 99, 235, 0.16);
  transform: translateY(-1px);
}

.category-filter-chip.is-dimmed {
  opacity: 0.65;
  background: var(--bg-ffffff);
}

.category-filter-chip.is-dimmed:hover {
  opacity: 1;
}

.chip-status-icon {
  width: 14px;
  height: 14px;
  border-radius: 50%;
  background: #2563eb;
  color: #ffffff;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.chip-color-dot {
  width: 9px;
  height: 9px;
  border-radius: 50%;
  flex-shrink: 0;
  box-shadow: 0 0 0 1px rgba(0, 0, 0, 0.06);
}

.chip-label {
  line-height: 1.3;
  text-align: left;
}

.chip-count-pill {
  font-size: 0.75rem;
  font-weight: 600;
  padding: 1px 7px;
  border-radius: 12px;
  background: var(--bg-f1f5f9);
  color: var(--tx-64748b);
  flex-shrink: 0;
  transition: all 0.2s ease;
}

.category-filter-chip.active-chip .chip-count-pill {
  background: var(--bg-dbeafe);
  color: var(--tx-1e40af);
}
/* Project Area Details Banner */
.project-area-banner {
  background: var(--bg-f8fafc);
  border: 1px solid var(--bd-e2e8f0);
  border-radius: 14px;
  padding: 16px 20px;
  box-shadow: 0 1px 3px rgba(15, 23, 42, 0.03);
  transition: all 0.25s ease;
}

.area-banner-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  flex-wrap: wrap;
  gap: 12px;
  margin-bottom: 14px;
  padding-bottom: 12px;
  border-bottom: 1px solid var(--bd-e2e8f0);
}

.banner-title-group {
  display: flex;
  align-items: center;
  gap: 10px;
  flex-wrap: wrap;
}

.banner-badge {
  width: 12px;
  height: 12px;
  border-radius: 4px;
  flex-shrink: 0;
}

.banner-title {
  font-size: 1rem;
  font-weight: 700;
  color: var(--tx-0f172a);
}

.banner-category-tag {
  font-size: 0.78rem;
  font-weight: 500;
  color: var(--tx-64748b);
  background: var(--bg-ffffff);
  padding: 2px 10px;
  border-radius: 6px;
  border: 1px solid var(--bd-e2e8f0);
}

.btn-select-project-action {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  padding: 6px 14px;
  border-radius: 8px;
  background: #2563eb;
  color: #ffffff;
  border: none;
  font-size: 0.82rem;
  font-weight: 600;
  cursor: pointer;
  transition: background 0.15s ease, transform 0.1s ease;
}

.btn-select-project-action:hover {
  background: #1d4ed8;
  transform: translateY(-1px);
}

.area-metrics-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(180px, 1fr));
  gap: 12px;
}

.metric-card {
  background: var(--bg-ffffff);
  border: 1px solid var(--bd-eef2f6);
  border-radius: 10px;
  padding: 12px 14px;
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.metric-label {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 0.75rem;
  color: var(--tx-64748b);
  font-weight: 500;
}

.metric-label svg {
  color: var(--tx-2563eb);
}

.metric-value {
  font-size: 1.05rem;
  font-weight: 700;
  color: var(--tx-1e293b);
}

.metric-value.highlight {
  color: var(--tx-2563eb);
}

.metric-value .unit {
  font-size: 0.78rem;
  font-weight: 500;
  color: var(--tx-64748b);
}

.metric-value.text-danger {
  color: var(--tx-dc2626);
}

/* ===== Drill-down Panel (inside map) ===== */
.stall-mode-panel {
  position: absolute;
  top: 14px;
  left: 14px;
  bottom: 14px;
  width: 320px;
  z-index: 12;
  display: flex;
  flex-direction: column;
  gap: 12px;
  padding: 16px;
  overflow-y: auto;
  background: var(--bg-rgba-255-255-255-0p96);
  backdrop-filter: blur(10px);
  border: 1px solid var(--bd-e2e8f0);
  border-radius: 14px;
  box-shadow: 0 10px 30px rgba(15, 23, 42, 0.12);
}

.smp-back-section {
  width: 100%;
  padding-bottom: 13px;
  border-bottom: 1px solid var(--bd-e8eef6);
}

.smp-back-btn {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  align-self: flex-start;
  padding: 6px 12px;
  border-radius: 8px;
  border: 1px solid #2563eb;
  background: #2563eb;
  color: #ffffff;
  font-family: inherit;
  font-size: 0.875rem;
  font-weight: 700;
  cursor: pointer;
  box-shadow: 0 4px 10px rgba(37, 99, 235, 0.22);
  transition: background 0.15s ease, box-shadow 0.15s ease, transform 0.15s ease;
}

.smp-back-btn:hover {
  background: #1d4ed8;
  box-shadow: 0 6px 14px rgba(37, 99, 235, 0.3);
  transform: translateX(-2px);
}

.smp-title-row {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-top: 2px;
}

.smp-dot {
  width: 12px;
  height: 12px;
  border-radius: 4px;
  flex-shrink: 0;
}

.smp-title {
  margin: 0;
  font-size: 1.25rem;
  font-weight: 700;
  color: var(--tx-0f172a);
}

.smp-category {
  margin-top: -6px;
  font-size: 0.9375rem;
  color: var(--tx-64748b);
  line-height: 1.4;
}

/* Multi-floor Selector inside Stall Mode Panel */
.smp-floor-section {
  display: flex;
  flex-direction: column;
  gap: 8px;
  background: var(--bg-f8fafc);
  border: 1px solid var(--bd-e2e8f0);
  border-radius: 12px;
  padding: 10px 12px;
}

.smp-floor-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.smp-floor-title {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  font-size: 0.875rem;
  font-weight: 700;
  color: var(--tx-334155);
}

.smp-floor-count {
  font-size: 0.8125rem;
  font-weight: 600;
  color: var(--tx-2563eb);
  background: var(--bg-eff6ff);
  border: 1px solid var(--bd-bfdbfe);
  padding: 1px 7px;
  border-radius: 9999px;
}

.smp-floor-segmented-row {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(0, 1fr));
  gap: 6px;
  margin-top: 4px;
}

.smp-floor-segmented-btn {
  display: flex;
  align-items: center;
  gap: 7px;
  padding: 6px 9px;
  border-radius: 9px;
  border: 1.5px solid var(--bd-e2e8f0);
  background: var(--bg-ffffff);
  cursor: pointer;
  font-family: inherit;
  transition: all 0.18s cubic-bezier(0.4, 0, 0.2, 1);
  text-align: left;
}

.smp-floor-segmented-btn:hover {
  border-color: var(--bd-cbd5e1);
  background: var(--bg-f8fafc);
}

.smp-floor-segmented-btn.active {
  background: var(--bg-eff6ff);
  border-color: #2563eb;
  box-shadow: 0 2px 6px rgba(37, 99, 235, 0.12);
}

.smp-seg-badge {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 24px;
  height: 22px;
  border-radius: 6px;
  background: var(--bg-f1f5f9);
  color: var(--tx-334155);
  font-size: 0.8125rem;
  font-weight: 800;
  flex-shrink: 0;
  transition: all 0.18s ease;
}

.smp-floor-segmented-btn.active .smp-seg-badge {
  background: #2563eb;
  color: #ffffff;
}

.smp-seg-info {
  display: flex;
  flex-direction: column;
  line-height: 1.2;
  min-width: 0;
}

.smp-seg-name {
  font-size: 0.8125rem;
  font-weight: 700;
  color: var(--tx-1e293b);
  white-space: nowrap;
}

.smp-floor-segmented-btn.active .smp-seg-name {
  color: var(--tx-1d4ed8);
}

.smp-seg-avail {
  font-size: 0.75rem;
  font-weight: 600;
  color: var(--tx-16a34a);
  white-space: nowrap;
}

.smp-floor-segmented-btn.active .smp-seg-avail {
  color: var(--tx-15803d);
  font-weight: 700;
}

/* Three or more floors share the narrow side panel: stack each tab's content
   so its label never competes horizontally with the floor badge. */
.smp-floor-segmented-row.has-many-floors {
  gap: 4px;
}

.has-many-floors .smp-floor-segmented-btn {
  flex-direction: column;
  justify-content: center;
  gap: 3px;
  min-width: 0;
  padding: 7px 3px 6px;
  text-align: center;
}

.has-many-floors .smp-seg-info {
  align-items: center;
}

.has-many-floors .smp-seg-name {
  display: none;
}

.has-many-floors .smp-seg-avail {
  font-size: 0.75rem;
}

/* Large buildings keep floor selection usable without making the panel taller
   than the map: two readable columns with an internal scroll area. */
.smp-floor-segmented-row.has-five-or-more-floors {
  grid-template-columns: repeat(2, minmax(0, 1fr));
  max-height: 178px;
  overflow-y: auto;
  padding-right: 3px;
  scrollbar-width: thin;
  scrollbar-color: #cbd5e1 transparent;
}

.has-five-or-more-floors .smp-floor-segmented-btn {
  flex-direction: row;
  justify-content: flex-start;
  padding: 7px 9px;
  text-align: left;
}

.has-five-or-more-floors .smp-seg-info {
  align-items: flex-start;
}

.has-five-or-more-floors .smp-seg-avail {
  font-size: 0.75rem;
}

.smp-floor-segmented-btn.is-floor-full .smp-seg-avail.is-full {
  color: var(--tx-dc2626);
  font-weight: 700;
}

.smp-floor-segmented-btn.is-floor-queue .smp-seg-avail.is-queue {
  color: var(--tx-d97706);
  font-weight: 700;
}

.smp-floor-segmented-btn.is-floor-full:not(.active) {
  opacity: 0.8;
  background: var(--bg-fafbfc);
}

.smp-floor-segmented-btn.is-floor-queue:not(.active) {
  background: var(--bg-fffdf5);
}

.smp-floor-segmented-btn.is-floor-full.active {
  border-color: var(--bd-fca5a5);
  background: var(--bg-fef2f2);
}

.smp-floor-segmented-btn.is-floor-full.active .smp-seg-badge {
  background: #dc2626;
  color: #ffffff;
}

.smp-floor-segmented-btn.is-floor-full.active .smp-seg-name {
  color: var(--tx-b91c1c);
}

.smp-floor-segmented-btn.is-floor-queue.active {
  border-color: #fcd34d;
  background: var(--bg-fffbeb);
}

.smp-floor-segmented-btn.is-floor-queue.active .smp-seg-badge {
  background: #f59e0b;
  color: #ffffff;
}

.smp-floor-segmented-btn.is-floor-queue.active .smp-seg-name {
  color: var(--tx-b45309);
}

.smp-floor-full-banner {
  background: linear-gradient(135deg, var(--bg-fef2f2) 0%, var(--bg-fff7ed) 100%);
  border: 1px solid var(--bd-fecaca);
  border-radius: 12px;
  padding: 11px 13px;
  display: flex;
  flex-direction: column;
  gap: 6px;
  box-shadow: 0 2px 8px rgba(220, 38, 38, 0.06);
}

.sffb-head {
  display: flex;
  align-items: center;
  gap: 6px;
  color: var(--tx-b91c1c);
  font-size: 0.84rem;
  font-weight: 700;
}

.sffb-head svg {
  flex-shrink: 0;
  color: var(--tx-dc2626);
}

.sffb-desc {
  margin: 0;
  font-size: 0.8125rem;
  color: var(--tx-7f1d1d);
  line-height: 1.45;
}

.sffb-switch-btn {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: 6px;
  width: 100%;
  padding: 7px 12px;
  margin-top: 2px;
  background: #2563eb;
  color: #ffffff;
  border: none;
  border-radius: 8px;
  font-family: inherit;
  font-size: 0.78rem;
  font-weight: 700;
  cursor: pointer;
  transition: all 0.15s ease;
  box-shadow: 0 2px 8px rgba(37, 99, 235, 0.25);
}

.sffb-switch-btn:hover {
  background: #1d4ed8;
  transform: translateY(-1px);
  box-shadow: 0 4px 12px rgba(37, 99, 235, 0.35);
}

.smp-floor-queue-banner {
  background: linear-gradient(135deg, var(--bg-fffbeb) 0%, var(--bg-fef3c7) 100%);
  border: 1px solid var(--bd-fde68a);
  border-radius: 12px;
  padding: 11px 13px;
  display: flex;
  flex-direction: column;
  gap: 6px;
  box-shadow: 0 2px 8px rgba(217, 119, 6, 0.08);
}

.sfqb-head {
  display: flex;
  align-items: center;
  gap: 6px;
  color: var(--tx-b45309);
  font-size: 0.84rem;
  font-weight: 700;
}

.sfqb-head svg {
  flex-shrink: 0;
  color: var(--tx-d97706);
}

.sfqb-desc {
  margin: 0;
  font-size: 0.8125rem;
  color: var(--tx-92400e);
  line-height: 1.45;
}

.sfqb-switch-btn {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: 6px;
  width: 100%;
  padding: 7px 12px;
  margin-top: 2px;
  background: #2563eb;
  color: #ffffff;
  border: none;
  border-radius: 8px;
  font-family: inherit;
  font-size: 0.78rem;
  font-weight: 700;
  cursor: pointer;
  transition: all 0.15s ease;
  box-shadow: 0 2px 8px rgba(37, 99, 235, 0.25);
}

.sfqb-switch-btn:hover {
  background: #1d4ed8;
  transform: translateY(-1px);
  box-shadow: 0 4px 12px rgba(37, 99, 235, 0.35);
}

/* Compact project status and searchable stall list */
.smp-stats-compact {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 5px;
  padding: 8px;
  border: 1px solid var(--bd-e2e8f0);
  border-radius: 10px;
  background: var(--bg-f8fafc);
}

.smp-sc-item {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  min-width: 0;
  color: var(--tx-64748b);
  font-size: 0.75rem;
  white-space: nowrap;
}

.smp-sc-item strong {
  color: var(--tx-1e293b);
  font-size: 0.8125rem;
}

.smp-sc-dot {
  width: 7px;
  height: 7px;
  flex: 0 0 7px;
  border-radius: 9999px;
}

.smp-search-section {
  display: flex;
  flex-direction: column;
  gap: 7px;
}

.smp-search-box {
  position: relative;
  display: flex;
  align-items: center;
}

.smp-search-icon {
  position: absolute;
  left: 10px;
  color: var(--tx-94a3b8);
  pointer-events: none;
}

.smp-search-input {
  width: 100%;
  min-width: 0;
  height: 34px;
  padding: 0 31px 0 31px;
  border: 1px solid var(--bd-cbd5e1);
  border-radius: 8px;
  outline: none;
  background: var(--bg-ffffff);
  color: var(--tx-1e293b);
  font: inherit;
  font-size: 0.78rem;
  transition: border-color 0.15s ease, box-shadow 0.15s ease;
}

.smp-search-input::placeholder { color: var(--tx-94a3b8); }

.smp-search-input:focus {
  border-color: #3b82f6;
  box-shadow: 0 0 0 3px rgba(59, 130, 246, 0.13);
}

.smp-search-clear {
  position: absolute;
  right: 6px;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 22px;
  height: 22px;
  padding: 0;
  border: 0;
  border-radius: 5px;
  background: transparent;
  color: var(--tx-94a3b8);
  cursor: pointer;
}

.smp-search-clear:hover { background: var(--bg-f1f5f9); color: var(--tx-475569); }

.smp-filter-tabs {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 4px;
}

.smp-filter-tab {
  min-width: 0;
  padding: 5px 3px;
  border: 1px solid var(--bd-e2e8f0);
  border-radius: 6px;
  background: var(--bg-ffffff);
  color: var(--tx-64748b);
  font: inherit;
  font-size: 0.75rem;
  font-weight: 600;
  line-height: 1.15;
  white-space: nowrap;
  cursor: pointer;
}

.smp-filter-tab:hover { border-color: var(--bd-93c5fd); color: var(--tx-2563eb); }
.smp-filter-tab.active { border-color: #2563eb; background: var(--bg-eff6ff); color: var(--tx-1d4ed8); }
.smp-filter-tab.is-avail.active { border-color: #22c55e; background: var(--bg-f0fdf4); color: var(--tx-15803d); }
.smp-filter-tab.is-res.active { border-color: #f59e0b; background: var(--bg-fffbeb); color: var(--tx-b45309); }

.smp-stall-list {
  display: flex;
  flex-direction: column;
  gap: 7px;
  min-height: 0;
  max-height: 290px;
  overflow-y: auto;
  padding: 2px 3px 2px 0;
  scrollbar-width: thin;
  scrollbar-color: #cbd5e1 transparent;
}

.smp-stall-card {
  display: flex;
  flex-direction: column;
  gap: 7px;
  padding: 9px 10px;
  border: 1px solid var(--bd-e2e8f0);
  border-radius: 9px;
  background: var(--bg-ffffff);
  cursor: pointer;
  transition: border-color 0.15s ease, box-shadow 0.15s ease, transform 0.15s ease;
}

.smp-stall-card:hover { border-color: var(--bd-93c5fd); box-shadow: 0 3px 9px rgba(37, 99, 235, 0.1); transform: translateY(-1px); }
.smp-stall-card.is-selected { border-color: #2563eb; background: var(--bg-eff6ff); box-shadow: 0 3px 10px rgba(37, 99, 235, 0.16); }
.smp-stall-card.is-occupied { background: var(--bg-f8fafc); cursor: not-allowed; opacity: 0.76; }
.projects-master-container.is-management .smp-stall-card.is-occupied { cursor: pointer; opacity: 1; }

.smp-card-top, .smp-card-bottom, .smp-card-code-grp, .smp-card-specs, .smp-card-action {
  display: flex;
  align-items: center;
}

.smp-card-top, .smp-card-bottom { justify-content: space-between; gap: 8px; }
.smp-card-code-grp { min-width: 0; gap: 5px; }
.smp-card-code { color: var(--tx-0f172a); font-size: 0.82rem; white-space: nowrap; }
.smp-card-floor-badge { padding: 2px 5px; border-radius: 4px; background: var(--bg-f1f5f9); color: var(--tx-475569); font-size: 0.75rem; font-weight: 700; white-space: nowrap; }

.smp-card-status-badge, .smp-action-pill {
  display: inline-flex;
  align-items: center;
  gap: 3px;
  padding: 3px 6px;
  border-radius: 9999px;
  font-size: 0.75rem;
  font-weight: 700;
  white-space: nowrap;
}

.badge-available { background: var(--bg-dcfce7); color: var(--tx-15803d); }
.badge-reserved { background: var(--bg-fef3c7); color: var(--tx-b45309); }
.badge-occupied { background: var(--bg-e2e8f0); color: var(--tx-475569); }
.badge-selected { background: #2563eb; color: #ffffff; }

.smp-card-specs { min-width: 0; gap: 5px; color: var(--tx-64748b); font-size: 0.75rem; }
.smp-spec-sep { color: var(--tx-cbd5e1); }
.smp-spec-rent { color: var(--tx-1e293b); font-weight: 700; white-space: nowrap; }
.smp-spec-rent small { color: var(--tx-64748b); font-size: 0.75rem; font-weight: 500; }
.smp-card-action { margin-left: auto; }
.smp-action-pill.is-selected { background: var(--bg-dbeafe); color: var(--tx-1d4ed8); }
.smp-action-pill.is-avail { background: var(--bg-dcfce7); color: var(--tx-15803d); }
.smp-action-pill.is-wait { background: var(--bg-fef3c7); color: var(--tx-b45309); }
.smp-action-pill.is-occ { background: var(--bg-e2e8f0); color: var(--tx-64748b); }

.smp-empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 7px;
  padding: 22px 10px;
  color: var(--tx-64748b);
  font-size: 0.75rem;
  text-align: center;
}

.smp-reopen-floor-btn {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  font-size: 0.75rem;
  font-weight: 600;
  color: var(--tx-2563eb);
  background: var(--bg-eff6ff);
  border: 1px solid var(--bd-bfdbfe);
  padding: 3px 8px;
  border-radius: 6px;
  cursor: pointer;
  transition: all 0.16s ease;
  font-family: inherit;
}

.smp-reopen-floor-btn:hover {
  background: var(--bg-dbeafe);
  color: var(--tx-1d4ed8);
  transform: translateY(-1px);
}

/* ========================================================== */
/* FLOATING ACTIVE FLOOR BADGE ON MAP CANVAS                 */
/* ========================================================== */
.map-floating-active-floor-badge {
  position: absolute;
  top: 14px;
  left: 170px;
  z-index: 20;
  display: inline-flex;
  align-items: center;
  gap: 8px;
  background: var(--bg-rgba-255-255-255-0p96);
  backdrop-filter: blur(10px);
  border: 1px solid var(--bd-bfdbfe);
  border-radius: 9999px;
  padding: 5px 8px 5px 12px;
  box-shadow: 0 4px 14px rgba(37, 99, 235, 0.12);
  cursor: pointer;
  transition: all 0.2s cubic-bezier(0.16, 1, 0.3, 1);
}

.map-floating-active-floor-badge:hover {
  background: var(--bg-ffffff);
  border-color: #3b82f6;
  box-shadow: 0 6px 18px rgba(37, 99, 235, 0.2);
  transform: translateY(-1px);
}

.mf-floor-tag {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  font-size: 0.78rem;
  color: var(--tx-1e293b);
}

.mf-floor-tag svg {
  color: var(--tx-2563eb);
  flex-shrink: 0;
}

.mf-floor-tag strong {
  color: var(--tx-2563eb);
  font-weight: 800;
}

.mf-change-floor-btn {
  display: inline-flex;
  align-items: center;
  gap: 3px;
  font-size: 0.75rem;
  font-weight: 700;
  color: var(--tx-2563eb);
  background: var(--bg-eff6ff);
  border: 1px solid var(--bd-dbeafe);
  padding: 3px 9px;
  border-radius: 9999px;
  cursor: pointer;
  pointer-events: none; /* Whole capsule is clickable */
}

/* ========================================================== */
/* FLOOR SELECTION PERSPECTIVE (MODAL / OVERLAY)             */
/* ========================================================== */
.floor-selector-modal-overlay {
  position: absolute;
  inset: 0;
  z-index: 45;
  background: rgba(15, 23, 42, 0.45);
  backdrop-filter: blur(6px);
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 24px;
}

.floor-selector-panel {
  width: 100%;
  max-width: 620px;
  max-height: min(90vh, 90dvh);
  background: var(--bg-ffffff);
  border-radius: 20px;
  border: 1px solid var(--bd-rgba-226-232-240-0p9);
  box-shadow: 0 24px 48px -12px rgba(15, 23, 42, 0.25), 0 0 0 1px rgba(0, 0, 0, 0.04);
  padding: 0;
  display: flex;
  flex-direction: column;
  gap: 0;
  overflow: hidden;
  animation: fscScaleUp 0.26s cubic-bezier(0.16, 1, 0.3, 1);
}

@keyframes fscScaleUp {
  from {
    opacity: 0;
    transform: scale(0.96) translateY(8px);
  }
  to {
    opacity: 1;
    transform: scale(1) translateY(0);
  }
}

.fsc-top-bar {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 16px;
}

.fsc-top-left {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.fsc-badge-row {
  display: flex;
  align-items: center;
  gap: 14px;
  margin-bottom: 10px;
  flex-wrap: wrap;
}

.fsc-building-tag {
  display: inline-flex;
  align-items: center;
  gap: 5px;
  font-size: 0.75rem;
  font-weight: 700;
  padding: 3px 10px;
  border-radius: 9999px;
  border: 1px solid;
}

.fsc-proj-code {
  font-size: 0.75rem;
  font-weight: 600;
  color: var(--tx-64748b);
  background: var(--bg-f1f5f9);
  border: 1px solid var(--bd-e2e8f0);
  padding: 3px 9px;
  border-radius: 9999px;
}

.fsc-title {
  font-size: 1.25rem;
  font-weight: 800;
  color: var(--tx-0f172a);
  margin: 2px 0 0;
  line-height: 1.3;
}

.fsc-desc {
  font-size: 0.84rem;
  color: var(--tx-64748b);
  margin: 0;
}

.fsc-close-btn {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 36px;
  height: 36px;
  border-radius: 10px;
  border: 1px solid var(--bd-e2e8f0);
  background: var(--bg-ffffff);
  color: var(--tx-64748b);
  cursor: pointer;
  transition: all 0.16s ease;
  flex-shrink: 0;
}

.fsc-close-btn:hover {
  background: var(--bg-fee2e2);
  color: var(--tx-dc2626);
  border-color: var(--bd-fca5a5);
  transform: rotate(90deg);
}

/* Guidance Card */
.fsc-guidance-card {
  display: flex;
  align-items: center;
  gap: 12px;
  background: linear-gradient(135deg, var(--bg-eff6ff) 0%, var(--bg-f8fafc) 100%);
  border: 1px solid var(--bd-bfdbfe);
  border-radius: 14px;
  padding: 12px 16px;
}

.fsc-guidance-icon-box {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 38px;
  height: 38px;
  border-radius: 10px;
  background: #2563eb;
  color: #ffffff;
  flex-shrink: 0;
  box-shadow: 0 4px 10px rgba(37, 99, 235, 0.25);
}

.fsc-guidance-content {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.fsc-step-label {
  font-size: 0.78rem;
  font-weight: 800;
  color: var(--tx-1d4ed8);
  letter-spacing: 0.02em;
}

.fsc-step-hint {
  font-size: 0.8rem;
  color: var(--tx-475569);
  line-height: 1.4;
}

/* Floor Stack */
.fsc-floors-stack {
  flex: 1 1 auto;
  min-height: 0;
  display: flex;
  flex-direction: column;
  gap: 10px;
  overflow-y: auto;
  overscroll-behavior: contain;
  padding: 20px 28px 24px;
  box-sizing: border-box;
}

.fsc-floor-item {
  display: flex;
  align-items: center;
  gap: 14px;
  width: 100%;
  padding: 14px 16px;
  background: var(--bg-ffffff);
  border: 1.5px solid var(--bd-e2e8f0);
  border-radius: 14px;
  cursor: pointer;
  font-family: inherit;
  text-align: left;
  transition: all 0.2s cubic-bezier(0.16, 1, 0.3, 1);
  position: relative;
}

.fsc-floor-item:hover {
  border-color: #3b82f6;
  background: var(--bg-f8fafc);
  box-shadow: 0 6px 18px rgba(37, 99, 235, 0.1);
  transform: translateY(-2px);
}

.fsc-floor-item.is-active-preview {
  border-color: #2563eb;
  background: var(--bg-eff6ff);
  box-shadow: 0 6px 20px rgba(37, 99, 235, 0.14);
}

.fsc-pillar {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  min-width: 66px;
  width: auto;
  padding: 0 6px;
  height: 58px;
  border-radius: 12px;
  background: var(--bg-f1f5f9);
  border: 1.5px solid var(--bd-e2e8f0);
  flex-shrink: 0;
  transition: all 0.2s ease;
}

.fsc-pillar-level {
  font-size: 0.95rem;
  font-weight: 800;
  color: var(--tx-1e293b);
  line-height: 1.2;
  white-space: nowrap;
}

.fsc-pillar-badge {
  font-size: 0.75rem;
  font-weight: 700;
  color: var(--tx-64748b);
  margin-top: 3px;
}

.fsc-floor-item.is-active-preview .fsc-pillar,
.fsc-floor-item:hover .fsc-pillar {
  background: #2563eb;
  border-color: #1d4ed8;
}

.fsc-floor-item.is-active-preview .fsc-pillar .fsc-pillar-level,
.fsc-floor-item.is-active-preview .fsc-pillar .fsc-pillar-badge,
.fsc-floor-item:hover .fsc-pillar .fsc-pillar-level,
.fsc-floor-item:hover .fsc-pillar .fsc-pillar-badge {
  color: #ffffff;
}

.fsc-floor-details {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 6px;
  min-width: 0;
}

.fsc-floor-headline {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 10px;
}

.fsc-floor-title-text {
  font-size: 0.95rem;
  font-weight: 700;
  color: var(--tx-0f172a);
}

.fsc-avail-chip {
  display: inline-flex;
  align-items: center;
  gap: 5px;
  font-size: 0.75rem;
  font-weight: 700;
  padding: 3px 9px;
  border-radius: 9999px;
  background: var(--bg-dcfce7);
  color: var(--tx-15803d);
  border: 1px solid var(--bd-bbf7d0);
  flex-shrink: 0;
}

.fsc-avail-chip.is-none {
  background: var(--bg-f1f5f9);
  color: var(--tx-94a3b8);
  border-color: var(--bd-e2e8f0);
}

.fsc-avail-chip.is-full {
  background: var(--bg-fee2e2);
  color: var(--tx-dc2626);
  border-color: var(--bd-fca5a5);
}

.fsc-avail-chip.is-queue {
  background: var(--bg-fef3c7);
  color: var(--tx-d97706);
  border-color: var(--bd-fde68a);
}

.fsc-floor-item.is-floor-full {
  background: var(--bg-fafbfc);
  opacity: 0.88;
}

.fsc-floor-item.is-floor-queue {
  background: var(--bg-ffffff);
}

.fsc-floor-item.is-floor-full:hover {
  opacity: 1;
  background: var(--bg-fffafa);
  border-color: var(--bd-fecaca);
  box-shadow: 0 6px 18px rgba(220, 38, 38, 0.08);
}

.fsc-pillar.is-full {
  background: var(--bg-f1f5f9);
  border-color: var(--bd-cbd5e1);
}

.fsc-pillar.is-full .fsc-pillar-level,
.fsc-pillar.is-full .fsc-pillar-badge {
  color: var(--tx-64748b);
}

/* Do not apply the blue "available/select" hover treatment to a full floor. */
.fsc-floor-item.is-floor-full:hover .fsc-pillar {
  background: var(--bg-f1f5f9);
  border-color: var(--bd-cbd5e1);
}

.fsc-floor-item.is-floor-full:hover .fsc-pillar .fsc-pillar-level,
.fsc-floor-item.is-floor-full:hover .fsc-pillar .fsc-pillar-badge {
  color: var(--tx-64748b);
}

.fsc-meta-pill.is-full-tag {
  background: var(--bg-fee2e2);
  color: var(--tx-b91c1c);
  font-weight: 700;
  border-color: var(--bd-fecaca);
}

.fsc-select-action-btn.is-full-action {
  background: var(--bg-f8fafc);
  color: var(--tx-64748b);
  border: 1px solid var(--bd-cbd5e1);
  box-shadow: none;
}

.fsc-select-action-btn.is-queue-action {
  background: var(--bg-fffbeb);
  color: var(--tx-d97706);
  border: 1px solid var(--bd-fde68a);
  box-shadow: none;
}

.fsc-floor-item:hover .fsc-select-action-btn.is-queue-action {
  background: var(--bg-fef3c7);
  color: var(--tx-b45309);
  border-color: #fcd34d;
  box-shadow: 0 2px 6px rgba(217, 119, 6, 0.2);
}

.fsc-floor-item:hover .fsc-select-action-btn.is-full-action {
  background: var(--bg-f1f5f9);
  color: var(--tx-475569);
  border-color: var(--bd-94a3b8);
  box-shadow: 0 2px 6px rgba(71, 85, 105, 0.12);
}

.fsc-chip-dot {
  width: 6px;
  height: 6px;
  border-radius: 50%;
  background: currentColor;
}

.fsc-floor-meta-row {
  display: flex;
  align-items: center;
  gap: 8px;
  flex-wrap: wrap;
}

.fsc-meta-pill {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  font-size: 0.75rem;
  color: var(--tx-64748b);
  background: var(--bg-f8fafc);
  padding: 3px 8px;
  border-radius: 6px;
  border: 1px solid var(--bd-f1f5f9);
}

.fsc-meta-pill svg {
  color: var(--tx-94a3b8);
}

.fsc-action-col {
  flex-shrink: 0;
}

.fsc-select-action-btn {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  padding: 8px 14px;
  border-radius: 10px;
  background: #2563eb;
  color: #ffffff;
  border: none;
  font-size: 0.8rem;
  font-weight: 700;
  cursor: pointer;
  transition: all 0.16s ease;
  box-shadow: 0 2px 6px rgba(37, 99, 235, 0.25);
  font-family: inherit;
}

.fsc-floor-item:hover .fsc-select-action-btn {
  background: #1d4ed8;
  box-shadow: 0 4px 10px rgba(37, 99, 235, 0.35);
  transform: translateX(2px);
}

/* Actions Bar */
.fsc-actions-bar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  padding-top: 10px;
  border-top: 1px solid var(--bd-f1f5f9);
}

.fsc-back-to-overview-btn {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  padding: 7px 12px;
  border-radius: 8px;
  border: 1px solid var(--bd-cbd5e1);
  background: var(--bg-ffffff);
  color: var(--tx-475569);
  font-size: 0.78rem;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.16s ease;
  font-family: inherit;
}

.fsc-back-to-overview-btn:hover {
  background: var(--bg-f1f5f9);
  color: var(--tx-1e293b);
  border-color: var(--bd-94a3b8);
}

.fsc-footer-tip {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  font-size: 0.75rem;
  color: var(--tx-64748b);
}

.fsc-footer-tip svg {
  color: var(--tx-2563eb);
  flex-shrink: 0;
}

/* Fade Transition */
.fsc-fade-enter-active,
.fsc-fade-leave-active {
  transition: opacity 0.22s ease;
}

.fsc-fade-enter-from,
.fsc-fade-leave-to {
  opacity: 0;
}

@media (max-width: 640px) {
  .floor-selector-panel {
    padding: 18px 16px;
  }
  .fsc-floor-item {
    flex-direction: column;
    align-items: stretch;
    gap: 10px;
  }
  .fsc-action-col {
    display: flex;
    justify-content: flex-end;
  }
  .map-floating-active-floor-badge {
    left: 14px;
    top: 54px;
  }
}

.smp-area-box {
  padding: 12px 14px;
  border-radius: 10px;
  background: linear-gradient(135deg, var(--bg-eff6ff) 0%, var(--bg-ffffff) 100%);
  border: 1px solid var(--bd-bfdbfe);
}

.smp-area-label {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 0.75rem;
  font-weight: 500;
  color: var(--tx-1e40af);
}

.smp-area-value {
  margin-top: 2px;
  font-size: 1.5rem;
  font-weight: 700;
  color: var(--tx-1d4ed8);
  letter-spacing: -0.01em;
}

.smp-area-value span {
  font-size: 0.85rem;
  font-weight: 500;
  color: var(--tx-64748b);
}

.smp-stats {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 8px;
}

.smp-stats.is-loading .smp-stat {
  cursor: wait;
}

.smp-stat {
  display: flex;
  flex-direction: column;
  gap: 2px;
  padding: 8px 10px;
  border-radius: 8px;
  background: var(--bg-f8fafc);
  border: 1px solid var(--bd-eef2f6);
  color: inherit;
  font-family: inherit;
  text-align: left;
  cursor: pointer;
  transition: border-color 0.16s ease, background 0.16s ease, box-shadow 0.16s ease, transform 0.16s ease;
}

.smp-stat:hover { border-color: var(--bd-bfdbfe); background: var(--bg-eff6ff); }
.smp-stat:focus-visible { outline: 2px solid var(--bd-93c5fd); outline-offset: 2px; }
.smp-stat.is-active { border-color: #2563eb; background: var(--bg-eff6ff); box-shadow: 0 0 0 2px rgba(37, 99, 235, 0.1); }
.smp-stat:active { transform: translateY(1px); }
.smp-stat:disabled {
  opacity: 1;
  pointer-events: none;
}

.smp-stat-skeleton {
  display: block;
  width: 58px;
  height: 25px;
  border-radius: 7px;
  background: linear-gradient(90deg, #e8edf4 25%, #f8fafc 50%, #e8edf4 75%);
  background-size: 200% 100%;
  animation: smp-stat-skeleton-loading 1.15s ease-in-out infinite;
}

@keyframes smp-stat-skeleton-loading {
  from { background-position: 200% 0; }
  to { background-position: -200% 0; }
}

.smp-stat-label {
  font-size: 0.8125rem;
  color: var(--tx-64748b);
}

.smp-stat strong {
  font-size: 1.25rem;
  font-weight: 700;
  color: var(--tx-1e293b);
}

.smp-stat.is-available strong { color: var(--tx-15803d); }
.smp-stat.is-reserved strong { color: var(--tx-d97706); }
.smp-stat.is-occupied strong { color: var(--tx-64748b); }

.smp-legend {
  display: flex;
  flex-direction: column;
  gap: 6px;
  padding: 10px 12px;
  background: var(--bg-ffffff);
  border-radius: 11px;
  border: 1px solid var(--bd-eef2f6);
  font-size: 0.875rem;
  color: var(--tx-334155);
}

.smp-legend-title {
  font-size: 0.875rem;
  font-weight: 600;
  color: var(--tx-64748b);
}

.smp-legend-items {
  display: flex;
  flex-wrap: wrap;
  gap: 8px 14px;
}

.smp-legend-items > div {
  display: inline-flex;
  align-items: center;
  gap: 6px;
}

.smp-scroll-content {
  display: contents;
}

/* Drag Handle Bar for Responsive Slide Up Sheet */
.smp-drag-handle-bar {
  display: none;
  width: 100%;
  padding: 8px 0 10px;
  cursor: grab;
  touch-action: none;
  justify-content: center;
  align-items: center;
  border-bottom: 1px solid var(--bd-f1f5f9);
  margin-bottom: 4px;
}
.smp-drag-handle-bar:active {
  cursor: grabbing;
}
.smp-drag-handle-pill {
  display: block;
  width: 44px;
  height: 5px;
  border-radius: 999px;
  background: var(--bg-cbd5e1);
  transition: background 0.2s ease, width 0.2s ease;
}
.smp-drag-handle-bar:hover .smp-drag-handle-pill {
  background: #94a3b8;
  width: 50px;
}
.smp-expand-toggle-btn {
  display: none;
}
.smp-header-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  width: 100%;
}

/* Occupancy Progress Bar Card */
.smp-occ-card {
  display: flex;
  flex-direction: column;
  gap: 8px;
  padding: 10px 14px;
  border: 1px solid var(--bd-e2e8f0);
  border-radius: 12px;
  background: var(--bg-f8fafc);
}
.smp-occ-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
}
.smp-occ-title-group {
  display: flex;
  flex-direction: column;
  gap: 1px;
}
.smp-occ-title {
  font-size: 0.9375rem;
  font-weight: 600;
  color: var(--tx-334155);
}
.smp-occ-sub {
  font-size: 0.8125rem;
  color: var(--tx-64748b);
}
.smp-occ-pct {
  font-size: 1.25rem;
  font-weight: 800;
  color: var(--tx-0f172a);
}
.smp-occ-track {
  width: 100%;
  height: 7px;
  border-radius: 999px;
  background: var(--bg-e2e8f0);
  display: flex;
  overflow: hidden;
}
.smp-occ-track .track-occupied {
  height: 100%;
  background: #2563eb;
  transition: width 0.3s ease;
}
.smp-occ-track .track-reserved {
  height: 100%;
  background: #f59e0b;
  transition: width 0.3s ease;
}
.smp-occ-legend {
  display: flex;
  flex-wrap: wrap;
  gap: 8px 12px;
  font-size: 0.8125rem;
  color: var(--tx-475569);
  font-weight: 500;
}
.smp-occ-legend span {
  display: inline-flex;
  align-items: center;
  gap: 5px;
}
.smp-occ-legend i {
  width: 8px;
  height: 8px;
  border-radius: 2px;
}
.dot-occ { background: #2563eb; }
.dot-res { background: #f59e0b; }
.dot-avail { background: #10b981; }

/* Project Info Grid (Area & Spec) */
.smp-project-info-grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 8px;
}
.smp-info-pill {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 9px 12px;
  border: 1px solid var(--bd-e2e8f0);
  border-radius: 11px;
  background: var(--bg-ffffff);
}
.smp-info-pill svg {
  color: var(--tx-2563eb);
  flex-shrink: 0;
}
.smp-info-text {
  display: flex;
  flex-direction: column;
}
.smp-info-text small {
  font-size: 0.8125rem;
  color: var(--tx-64748b);
  font-weight: 500;
  line-height: 1.2;
}
.smp-info-text strong {
  font-size: 1rem;
  font-weight: 700;
  color: var(--tx-0f172a);
  line-height: 1.3;
}
.smp-info-text strong em {
  font-size: 0.8125rem;
  font-weight: 500;
  font-style: normal;
  color: var(--tx-64748b);
}

.smp-sw {
  width: 12px;
  height: 12px;
  border-radius: 3px;
  flex-shrink: 0;
  border: 1px solid var(--bd-rgba-15-23-42-0p08);
}

.sw-available { background: #10b981; }
.sw-reserved { background: #f59e0b; }
.sw-occupied { background: #475569; }
.sw-selected { background: #2563eb; }

.smp-hint {
  display: flex;
  gap: 8px;
  align-items: flex-start;
  font-size: 0.8125rem;
  color: var(--tx-64748b);
  line-height: 1.45;
  padding: 9px 12px;
  background: var(--bg-f8fafc);
  border-radius: 10px;
  border: 1px solid var(--bd-eef2f6);
}

.smp-hint svg {
  flex-shrink: 0;
  margin-top: 2px;
  color: var(--tx-2563eb);
}

/* หน้าข้อมูลตลาดและฝั่งผู้เช่า: แสดงข้อมูลพื้นที่เป็นกล่องเดียวแบบสองแถว */
.smp-occ-card.is-map-colors { gap: 12px; padding: 14px 16px; }
.smp-occ-card.is-map-colors .track-occupied,
.smp-occ-card.is-map-colors .dot-occ { background: #475569; }
.smp-occ-card.is-map-colors .smp-occ-track { height: 8px; }
.smp-occ-card.is-map-colors .smp-occ-legend { justify-content: space-between; gap: 6px 14px; font-size: 0.8125rem; }
.smp-occ-card.is-map-colors .smp-occ-legend i { width: 8px; height: 8px; border-radius: 2px; }
.smp-project-info-grid.is-list {
  grid-template-columns: 1fr;
  gap: 0;
  border: 1px solid var(--bd-e2e8f0);
  border-radius: 12px;
  background: var(--bg-ffffff);
  overflow: hidden;
}
.smp-project-info-grid.is-list .smp-info-pill { min-height: 62px; padding: 11px 16px; border: 0; border-radius: 0; background: transparent; gap: 12px; box-sizing: border-box; }
.smp-project-info-grid.is-list .smp-info-pill + .smp-info-pill { border-top: 1px solid var(--bd-eef2f6); }
.smp-project-info-grid.is-list .smp-info-text { flex: 1 1 auto; min-width: 0; flex-direction: row; align-items: center; justify-content: space-between; gap: 12px; }
.smp-project-info-grid.is-list .smp-info-text small { flex: 1 1 auto; min-width: 0; font-size: 0.875rem; white-space: nowrap; }
.smp-project-info-grid.is-list .smp-info-text strong {
  display: flex;
  flex: 0 1 auto;
  flex-wrap: wrap;
  align-items: baseline;
  justify-content: flex-end;
  min-width: 0;
  max-width: 58%;
  column-gap: 4px;
  white-space: normal;
  text-align: right;
}
.smp-project-info-grid.is-list .smp-info-text strong em { flex: 0 0 auto; white-space: nowrap; }
.smp-project-info-grid.is-list .smp-info-text strong > span { flex: 0 0 auto; white-space: nowrap; }
.smp-project-info-grid.is-list .smp-info-text strong.smp-area-measure { flex-direction: column; flex-wrap: nowrap; align-items: flex-end; line-height: 1.15; }
.smp-hint.is-quiet { padding: 4px 2px 0; background: transparent; border: 0; border-radius: 0; font-size: 0.875rem; }

/* ช่วงวันที่ในตัวกรอง (ตั้งแต่ / ถึง) */
.ov-date-range { display: grid; grid-template-columns: 1fr 1fr; gap: 10px; margin-top: 8px; }
.ov-date-range label { display: flex; flex-direction: column; gap: 4px; min-width: 0; }
.ov-date-range label > span { color: var(--tx-64748b); font-size: 0.8125rem; font-weight: 500; }

.smp-slide-enter-active,
.smp-slide-leave-active {
  transition: opacity 0.25s ease, transform 0.25s ease;
}

.smp-slide-enter-from,
.smp-slide-leave-to {
  opacity: 0;
  transform: translateX(-16px);
}

@media (max-width: 768px) {
  .maplibre-overview-wrapper {
    height: 520px;
  }

  /* บนโหมด responsive ซ่อนปุ่มย้อนกลับและปุ่มลูกศร ใช้การรูดแผงลงเพื่อปิดกลับผังรวมแทน */
  .smp-desktop-only,
  .smp-header-row,
  .smp-back-section,
  .smp-expand-toggle-btn {
    display: none !important;
  }

  /* ขีดปรับขนาดตรึงอยู่ด้านบนสุดตลอดเวลา เป็นหัวกล่องขาวทึบ 100% ไม่โปร่งใส ไม่เห็นพื้นหลังทะลุ */
  .smp-drag-handle-bar {
    flex: 0 0 auto !important;
    display: flex !important;
    justify-content: center !important;
    align-items: center !important;
    width: 100% !important;
    padding: 12px 0 10px !important;
    margin: 0 !important;
    background: var(--bg-ffffff) !important;
    border-bottom: 1px solid var(--bd-f1f5f9) !important;
    border-radius: 20px 20px 0 0 !important;
    cursor: grab !important;
    touch-action: none !important;
    user-select: none !important;
    -webkit-user-select: none !important;
  }

  /* เมื่อเปิดโครงการในโหมด responsive ให้ซ่อน ov-rail ไม่ให้ลอยบังแผง slide up */
  .ov-rail.is-project-active {
    display: none !important;
    opacity: 0 !important;
    pointer-events: none !important;
  }

  .stall-mode-panel {
    position: fixed !important;
    top: auto !important;
    bottom: 0 !important;
    left: 0 !important;
    right: 0 !important;
    width: 100% !important;
    max-width: 100% !important;
    height: 50dvh !important;
    max-height: calc(100dvh - 58px) !important;
    border-radius: 20px 20px 0 0 !important;
    border: 1px solid var(--bd-e2e8f0) !important;
    border-bottom: 0 !important;
    box-shadow: 0 -8px 32px rgba(15, 23, 42, 0.18) !important;
    z-index: 200 !important;
    display: flex !important;
    flex-direction: column !important;
    background: var(--bg-ffffff) !important;
    padding: 0 !important;
    overflow: hidden !important;
    box-sizing: border-box !important;
    transition: height 0.28s cubic-bezier(0.16, 1, 0.3, 1), max-height 0.28s cubic-bezier(0.16, 1, 0.3, 1) !important;
  }

  .smp-scroll-content {
    flex: 1 1 auto !important;
    min-height: 0 !important;
    overflow-y: auto !important;
    -webkit-overflow-scrolling: touch !important;
    padding: 12px 16px 28px !important;
    display: flex !important;
    flex-direction: column !important;
    gap: 12px !important;
    background: var(--bg-ffffff) !important;
  }

  .stall-mode-panel.is-sheet-expanded {
    height: calc(100dvh - 58px) !important;
    max-height: calc(100dvh - 58px) !important;
  }

  .smp-slide-enter-from,
  .smp-slide-leave-to {
    opacity: 0;
    transform: translateY(100%) !important;
  }

  .smp-slide-enter-active,
  .smp-slide-leave-active {
    transition: transform 0.28s cubic-bezier(0.16, 1, 0.3, 1), opacity 0.22s ease !important;
  }

  .smp-back-section {
    padding-bottom: 0;
    border-bottom: none;
  }

  .smp-back-btn {
    padding: 5px 10px;
    font-size: 0.8125rem;
  }

  .smp-title {
    font-size: 1.0625rem;
  }

  .smp-category {
    font-size: 0.8125rem;
    margin-top: -2px;
  }

  /* Compact multi-floor row on mobile */
  .smp-floor-section {
    padding: 8px 10px;
    gap: 6px;
  }

  .smp-floor-header {
    font-size: 0.8125rem;
  }

  .smp-floor-segmented-row {
    display: flex;
    overflow-x: auto;
    flex-wrap: nowrap;
    gap: 6px;
    padding-bottom: 2px;
    -webkit-overflow-scrolling: touch;
  }

  .smp-floor-segmented-btn {
    flex: 0 0 auto !important;
    min-height: 40px !important;
    padding: 6px 12px !important;
    border-radius: 10px !important;
    display: inline-flex !important;
    align-items: center !important;
    gap: 7px !important;
  }

  .smp-floor-segmented-btn .smp-seg-badge {
    font-size: 0.8125rem !important;
    padding: 3px 7px !important;
    border-radius: 6px !important;
  }

  .smp-floor-segmented-btn .smp-seg-name {
    font-size: 0.875rem !important;
  }

  .smp-floor-segmented-btn .smp-seg-avail {
    font-size: 0.75rem !important;
  }

  .sffb-switch-btn,
  .sfqb-switch-btn {
    min-height: 46px !important;
    padding: 10px 16px !important;
    font-size: 0.875rem !important;
    border-radius: 10px !important;
    display: inline-flex !important;
    align-items: center !important;
    justify-content: center !important;
  }

  /* Show complete info on mobile */
  .smp-occ-card {
    display: flex !important;
  }

  .smp-project-info-grid {
    display: grid !important;
  }

  .smp-legend {
    display: flex !important;
  }

  .smp-hint {
    display: flex !important;
  }

  /* Compact 4-column stats on mobile */
  .smp-stats {
    grid-template-columns: repeat(4, 1fr);
    gap: 6px;
  }

  .smp-stat {
    padding: 8px 6px;
    text-align: center;
    border-radius: 10px;
  }

  .smp-stat-label {
    font-size: 0.75rem;
  }

  .smp-stat strong {
    font-size: 1.05rem;
  }
}

@media (max-width: 480px) {
  .map-floating-top-controls {
    top: 10px;
    right: 10px;
    gap: 6px;
    max-width: calc(100% - 20px);
    flex-wrap: wrap;
    justify-content: flex-end;
  }

  .market-floor-switcher-label {
    display: none;
  }

  .market-floor-tab {
    min-height: 38px;
    padding: 6px 12px;
    font-size: 0.8125rem;
    border-radius: 8px;
    display: inline-flex;
    align-items: center;
    justify-content: center;
  }

  .map-floating-style-switcher .map-mode-tab {
    min-height: 38px;
    padding: 6px 12px;
    font-size: 0.8125rem;
    border-radius: 8px;
    display: inline-flex;
    align-items: center;
    justify-content: center;
  }

  .map-floating-recenter-btn {
    min-height: 38px;
    padding: 6px 12px;
    font-size: 0.8125rem;
    border-radius: 8px;
    display: inline-flex;
    align-items: center;
    justify-content: center;
  }
}
</style>

<style>
/* Global MapLibre Marker Styling */
.maplibre-proj-node {
  cursor: pointer;
  transition: transform 0.2s ease, filter 0.2s ease;
  user-select: none;
}

/* Embedded Zone Number Stamp inside Building Footprint */
.maplibre-zone-number-badge {
  pointer-events: auto !important;
  cursor: pointer;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  background: rgba(15, 23, 42, 0.88);
  color: #ffffff;
  border: 1.5px solid var(--bd-rgba-255-255-255-0p95);
  border-radius: 999px;
  padding: 3px 8px;
  font-family: 'Bai Jamjuree', sans-serif;
  font-size: 12px;
  font-weight: 700;
  box-shadow: 0 2px 6px rgba(0, 0, 0, 0.35);
  /* ห้าม transition: all หรือ transform เพราะจะทำให้ Marker ดีเลย์/ลอยตามตอนลากแผนที่ */
  transition: background-color 0.15s ease, border-color 0.15s ease, box-shadow 0.15s ease, opacity 0.15s ease;
  user-select: none;
  white-space: nowrap;
}

.zone-badge-inner {
  display: inline-flex;
  align-items: center;
  gap: 5px;
  transition: transform 0.15s ease;
}

.zone-badge-inner .badge-num {
  font-weight: 800;
  color: #ffffff;
}

.zone-badge-inner .badge-floor {
  padding-left: 5px;
  border-left: 1px solid var(--bd-rgba-255-255-255-0p42);
  color: #bfdbfe;
  font-size: 12px;
  font-weight: 800;
}

.zone-badge-inner .badge-sep {
  opacity: 0.55;
  font-size: 0.75rem;
  color: var(--tx-94a3b8);
}

.zone-badge-inner .badge-area {
  font-size: 0.75rem;
  font-weight: 600;
  color: var(--tx-93c5fd);
}

.maplibre-zone-number-badge:hover .zone-badge-inner,
.maplibre-zone-number-badge.is-project-hovered .zone-badge-inner {
  transform: scale(1.15);
}

.maplibre-zone-number-badge:hover,
.maplibre-zone-number-badge.is-project-hovered {
  background: #2563eb;
  border-color: var(--bd-ffffff);
  box-shadow: 0 0 0 3px rgba(37, 99, 235, 0.4), 0 4px 12px rgba(0, 0, 0, 0.4);
}

.maplibre-zone-number-badge:hover .badge-area,
.maplibre-zone-number-badge.is-project-hovered .badge-area {
  color: #ffffff;
}

.maplibre-zone-number-badge.is-selected .zone-badge-inner {
  transform: scale(1.2);
}

.maplibre-zone-number-badge.is-selected {
  background: #2563eb;
  border-color: var(--bd-ffffff);
  box-shadow: 0 0 0 4px rgba(37, 99, 235, 0.45), 0 4px 14px rgba(0, 0, 0, 0.4);
}

.maplibre-zone-number-badge.is-selected .badge-area {
  color: #ffffff;
}

.maplibre-zone-number-badge.is-dimmed {
  opacity: 0.18;
}

/* At mid zoom, a small project number keeps the map readable without hiding polygons. */
.maplibre-zone-number-badge.is-compact {
  min-width: 22px;
  height: 22px;
  padding: 0 4px;
  border-radius: 50%;
  font-size: 12px;
}

.maplibre-zone-number-badge.is-compact.has-floor-context {
  width: auto;
  min-width: 38px;
  padding: 0 6px;
  border-radius: 999px;
}

.maplibre-zone-number-badge.is-floor-only {
  width: auto;
  min-width: 38px;
  padding: 3px 8px;
  border-radius: 999px;
}

.maplibre-zone-number-badge.is-floor-only .badge-floor {
  padding-left: 0;
  border-left: 0;
  color: #ffffff;
}

.maplibre-zone-number-badge.is-project-hovered .badge-floor,
.maplibre-zone-number-badge:hover .badge-floor {
  color: #ffffff;
}

/* Custom Popup Content (hover cards) */
.mlp-popup {
  z-index: 100 !important;
  pointer-events: none !important;
  filter: drop-shadow(0 0 1px rgba(15, 23, 42, 0.18))
          drop-shadow(0 14px 28px rgba(15, 23, 42, 0.22));
}
.mlp-popup:not(.is-interactive) .maplibregl-popup-content,
.mlp-popup:not(.is-interactive) .maplibregl-popup-tip {
  pointer-events: none !important;
}
.mlp-popup.is-interactive,
.mlp-popup.is-interactive .maplibregl-popup-content {
  pointer-events: auto !important;
}
.mlp-popup .maplibregl-popup-content {
  border-radius: 14px !important;
  border: none !important;
  box-shadow: none !important;
  padding: 0 !important;
  overflow: hidden;
  background: var(--bg-ffffff) !important;
}

/* Seamless Tip Connection - flush and fully unified with the card */
.mlp-popup .maplibregl-popup-tip {
  z-index: 2;
  position: relative;
}
.mlp-popup.maplibregl-popup-anchor-bottom .maplibregl-popup-tip,
.mlp-popup.maplibregl-popup-anchor-bottom-left .maplibregl-popup-tip,
.mlp-popup.maplibregl-popup-anchor-bottom-right .maplibregl-popup-tip {
  border-top-color: var(--bd-ffffff) !important;
  margin-top: -2px;
}
.mlp-popup.maplibregl-popup-anchor-top .maplibregl-popup-tip,
.mlp-popup.maplibregl-popup-anchor-top-left .maplibregl-popup-tip,
.mlp-popup.maplibregl-popup-anchor-top-right .maplibregl-popup-tip {
  border-bottom-color: var(--bd-ffffff) !important;
  margin-bottom: -2px;
}
.mlp-popup.maplibregl-popup-anchor-left .maplibregl-popup-tip {
  border-right-color: var(--bd-ffffff) !important;
  margin-right: -2px;
}
.mlp-popup.maplibregl-popup-anchor-right .maplibregl-popup-tip {
  border-left-color: var(--bd-ffffff) !important;
  margin-left: -2px;
}

.mlp {
  font-family: 'Bai Jamjuree', sans-serif;
  width: 360px;
  max-width: 100%;
}
.mlp-stall-tooltip { width: auto; min-width: 0; max-width: 220px; padding: 9px 14px; color: var(--tx-0f172a); font-size: 0.9375rem; font-weight: 700; line-height: 1.35; white-space: nowrap; }
.mlp-popup:has(.mlp-stall-tooltip) .maplibregl-popup-content { width: max-content; max-width: 240px; }
.mlp-accent {
  height: 4px;
  background: linear-gradient(90deg, var(--mlp-accent), color-mix(in srgb, var(--mlp-accent) 40%, var(--bg-ffffff)));
}
.mlp-photo-wrap {
  position: relative;
  width: 100%;
  height: 124px;
  overflow: hidden;
  background: var(--bg-f1f5f9);
}
.mlp-photo-img {
  width: 100%;
  height: 100%;
  object-fit: cover;
  display: block;
  transition: transform 0.3s ease;
}
.mlp:hover .mlp-photo-img {
  transform: scale(1.05);
}
.mlp-photo-overlay {
  position: absolute;
  top: 8px;
  left: 8px;
  right: 8px;
  display: flex;
  justify-content: flex-start;
  align-items: center;
  pointer-events: none;
}
.mlp-photo-badge {
  background: rgba(15, 23, 42, 0.72);
  backdrop-filter: blur(6px);
  color: #ffffff;
  font-size: 0.75rem;
  font-weight: 600;
  padding: 3px 8px;
  border-radius: 20px;
  display: inline-flex;
  align-items: center;
  gap: 4px;
}
.mlp-photo-badge svg {
  color: var(--tx-93c5fd);
}
.mlp-body { padding: 14px 16px 12px; }
.mlp-project-context { display: flex; align-items: center; gap: 6px; margin-bottom: 5px; overflow: hidden; color: var(--tx-64748b); font-size: .68rem; font-weight: 700; white-space: nowrap; }
.mlp-project-context span { overflow: hidden; text-overflow: ellipsis; }
.mlp-project-context i { flex: 0 0 auto; color: var(--tx-94a3b8); font-style: normal; }

.mlp-head {
  display: grid;
  grid-template-columns: minmax(0, 1fr) auto;
  align-items: start;
  gap: 8px;
}
.mlp-title {
  display: flex;
  align-items: center;
  gap: 8px;
  min-width: 0;
}
.mlp-stall-heading {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  min-width: 0;
}
.mlp-title strong {
  font-size: 1rem;
  font-weight: 700;
  color: var(--tx-0f172a);
  line-height: 1.35;
  overflow-wrap: anywhere;
}
.mlp-dot {
  width: 10px;
  height: 10px;
  border-radius: 3px;
  background: var(--mlp-accent);
  box-shadow: 0 0 0 3px color-mix(in srgb, var(--mlp-accent) 22%, transparent);
  flex-shrink: 0;
}
.mlp-pill {
  font-size: 0.75rem;
  font-weight: 600;
  padding: 3px 10px;
  border-radius: 999px;
  white-space: nowrap;
  flex-shrink: 0;
}
.mlp-pill.is-open { background: var(--bg-dcfce7); color: var(--tx-15803d); }
.mlp-pill.is-full { background: var(--bg-fee2e2); color: var(--tx-dc2626); }

.mlp-building-row {
  margin: 5px 0 0 18px;
}

.mlp-building-badge {
  display: inline-flex;
  align-items: center;
  font-size: 0.75rem;
  padding: 2px 7px;
  border-radius: 4px;
  background: var(--bg-eff6ff);
  color: var(--tx-2563eb);
  font-weight: 700;
  border: 1px solid var(--bd-bfdbfe);
}

.mlp-sub {
  margin-top: 4px;
  font-size: 0.76rem;
  color: var(--tx-64748b);
  line-height: 1.45;
}

.mlp-tiles {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 8px;
  margin-top: 12px;
}
.mlp-tile {
  background: var(--bg-f8fafc);
  border: 1px solid var(--bd-eef2f6);
  border-radius: 10px;
  padding: 8px 10px;
  min-width: 0;
}
.mlp-tile-label {
  display: flex;
  align-items: center;
  gap: 5px;
  font-size: 0.75rem;
  font-weight: 500;
  color: var(--tx-64748b);
  margin-bottom: 3px;
}
.mlp-tile-label svg { color: var(--tx-94a3b8); flex-shrink: 0; }
.mlp-tile-value {
  display: block;
  font-size: 0.95rem;
  font-weight: 700;
  color: var(--tx-1e293b);
  white-space: nowrap;
}
.mlp-tile-value small { font-size: 0.75rem; font-weight: 500; color: var(--tx-64748b); }
.mlp-tile-value.is-primary { color: var(--tx-2563eb); }
.mlp-tile-value.is-money { color: var(--tx-059669); }
.mlp-tile-note { display: block; font-size: 0.75rem; color: var(--tx-94a3b8); margin-top: 1px; }

.mlp-occ { margin-top: 12px; }
.mlp-occ-head {
  display: flex;
  justify-content: space-between;
  font-size: 0.75rem;
  color: var(--tx-64748b);
  margin-bottom: 6px;
}
.mlp-occ-pct { font-weight: 600; color: var(--tx-334155); }
.mlp-occ-bar {
  display: flex;
  height: 8px;
  border-radius: 999px;
  overflow: hidden;
  background: var(--bg-f1f5f9);
}
.mlp-occ-bar i { display: block; height: 100%; }
.mlp-occ-legend {
  display: flex;
  justify-content: space-between;
  margin-top: 7px;
  font-size: 0.75rem;
  color: var(--tx-64748b);
}
.mlp-occ-legend span { display: inline-flex; align-items: center; gap: 4px; }
.mlp-occ-legend i { width: 8px; height: 8px; border-radius: 2px; display: inline-block; }
.mlp-occ-legend b { color: var(--tx-334155); font-weight: 700; }

.mlp-hint {
  margin-top: 12px;
  padding-top: 9px;
  border-top: 1px dashed var(--bd-e2e8f0);
  font-size: 0.75rem;
  color: var(--tx-94a3b8);
  text-align: center;
}
.mlp-hint.is-reserved {
  color: var(--tx-d97706);
  font-weight: 500;
}
.mlp-facility-delete {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: 6px;
  width: 100%;
  height: 34px;
  margin-top: 10px;
  border: 1px solid var(--bd-fecaca);
  border-radius: 8px;
  background: var(--bg-ffffff);
  color: var(--tx-dc2626);
  font: 600 .8125rem/1 'Bai Jamjuree', sans-serif;
  cursor: pointer;
}
.mlp-facility-delete:hover { background: var(--bg-fef2f2); }
.mlp-facility-delete svg { width: 15px; height: 15px; }

/* Queue Info Box for Reserved Stalls */
.mlp-queue-box {
  display: flex;
  align-items: center;
  gap: 10px;
  margin-top: 10px;
  padding: 8px 12px;
  background: var(--bg-fffbeb);
  border: 1px solid var(--bd-fde68a);
  border-radius: 10px;
}
.mlp-queue-icon {
  flex-shrink: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  width: 28px;
  height: 28px;
  background: var(--bg-fef3c7);
  color: var(--tx-d97706);
  border-radius: 8px;
}
.mlp-queue-content {
  display: flex;
  flex-direction: column;
  gap: 2px;
  min-width: 0;
}
.mlp-queue-title {
  font-size: 0.76rem;
  font-weight: 700;
  color: var(--tx-92400e);
}
.mlp-queue-title b {
  color: var(--tx-d97706);
}
.mlp-queue-desc {
  font-size: 0.75rem;
  color: var(--tx-a16207);
  line-height: 1.35;
}
.mlp-queue-desc b {
  color: var(--tx-78350f);
  font-weight: 700;
}

/* Floating Map Style Switcher & Recenter Controls (inside map - Top Right) */
.map-floating-top-controls {
  position: absolute;
  top: 14px;
  right: 14px;
  z-index: 20;
  display: inline-flex;
  align-items: center;
  gap: 10px;
}

.market-floor-switcher {
  display: inline-flex;
  align-items: center;
  gap: 2px;
  padding: 3px;
  border: 1px solid var(--bd-rgba-191-219-254-0p95);
  border-radius: 9999px;
  background: var(--bg-rgba-239-246-255-0p96);
  backdrop-filter: blur(8px);
  box-shadow: 0 4px 14px rgba(37, 99, 235, 0.14);
}

.market-floor-switcher-label {
  display: inline-flex;
  align-items: center;
  gap: 5px;
  padding: 0 8px 0 7px;
  color: var(--tx-475569);
  font-size: 0.78rem;
  font-weight: 600;
  white-space: nowrap;
}

.market-floor-switcher-label svg {
  color: var(--tx-2563eb);
}

.market-floor-tab {
  border: none;
  border-radius: 9999px;
  padding: 6px 12px;
  background: transparent;
  color: var(--tx-64748b);
  font-family: inherit;
  font-size: 0.82rem;
  font-weight: 700;
  line-height: 1.2;
  cursor: pointer;
  transition: color 0.18s ease, background 0.18s ease, box-shadow 0.18s ease;
}

.market-floor-tab:hover {
  color: var(--tx-1d4ed8);
}

.market-floor-tab.active {
  color: #ffffff;
  background: #2563eb;
  box-shadow: 0 2px 8px rgba(37, 99, 235, 0.28);
}

.map-floating-style-switcher {
  display: inline-flex;
  align-items: center;
  background: var(--bg-rgba-255-255-255-0p95);
  backdrop-filter: blur(8px);
  padding: 3px;
  border-radius: 9999px;
  border: 1px solid var(--bd-rgba-226-232-240-0p9);
  box-shadow: 0 4px 14px rgba(15, 23, 42, 0.12);
}

.map-floating-style-switcher .map-mode-tab {
  padding: 6px 16px;
  border-radius: 9999px;
  border: none;
  background: transparent;
  font-family: inherit;
  font-size: 0.84rem;
  font-weight: 600;
  color: var(--tx-64748b);
  cursor: pointer;
  transition: all 0.18s ease;
  line-height: 1.2;
}

.map-floating-style-switcher .map-mode-tab:hover {
  color: var(--tx-1e293b);
}

.map-floating-style-switcher .map-mode-tab.active {
  background: var(--bg-ffffff);
  color: var(--tx-2563eb);
  font-weight: 700;
  box-shadow: 0 2px 8px rgba(37, 99, 235, 0.2);
}

.map-floating-recenter-btn {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  background: var(--bg-rgba-255-255-255-0p95);
  backdrop-filter: blur(8px);
  padding: 6px 14px;
  border-radius: 9999px;
  border: 1px solid var(--bd-rgba-226-232-240-0p9);
  box-shadow: 0 4px 14px rgba(15, 23, 42, 0.12);
  font-family: inherit;
  font-size: 0.84rem;
  font-weight: 600;
  color: var(--tx-334155);
  cursor: pointer;
  transition: all 0.18s ease;
  line-height: 1.2;
}

.map-floating-recenter-btn:hover {
  background: var(--bg-ffffff);
  color: var(--tx-2563eb);
  border-color: var(--bd-93c5fd);
  transform: translateY(-1px);
  box-shadow: 0 6px 16px rgba(37, 99, 235, 0.18);
}

.map-floating-recenter-btn:active {
  transform: translateY(0);
}

.map-floating-recenter-btn svg {
  color: var(--tx-2563eb);
  flex-shrink: 0;
}

.map-floating-floor-preview-pill {
  position: absolute;
  top: 14px;
  left: 14px;
  z-index: 20;
  display: inline-flex;
  align-items: center;
  gap: 8px;
  background: var(--bg-rgba-255-255-255-0p95);
  backdrop-filter: blur(8px);
  padding: 4px 10px 4px 12px;
  border-radius: 9999px;
  border: 1px solid var(--bd-rgba-226-232-240-0p9);
  box-shadow: 0 4px 14px rgba(15, 23, 42, 0.12);
}

.floor-preview-label {
  font-size: 0.78rem;
  font-weight: 700;
  color: var(--tx-64748b);
}

.floor-preview-tabs {
  display: inline-flex;
  align-items: center;
  gap: 4px;
}

.floor-preview-tab {
  padding: 4px 12px;
  border-radius: 9999px;
  border: none;
  background: transparent;
  color: var(--tx-475569);
  font-family: inherit;
  font-size: 0.8rem;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.15s ease;
}

.floor-preview-tab:hover {
  color: var(--tx-1e293b);
  background: var(--bg-f1f5f9);
}

.floor-preview-tab.active {
  background: #2563eb;
  color: #ffffff;
  font-weight: 700;
  box-shadow: 0 2px 6px rgba(37, 99, 235, 0.25);
}

/* Preview Mode Styles */
.projects-master-container.is-preview-mode {
  padding: 12px;
  gap: 12px;
  background: var(--bg-ffffff);
  border-radius: 14px;
}

.projects-master-container.is-preview-mode .maplibre-overview-wrapper {
  height: 540px;
}

@media (max-width: 768px) {
  .projects-master-container.is-preview-mode .maplibre-overview-wrapper {
    height: 440px;
  }
}

.preview-top-bar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  flex-wrap: wrap;
  gap: 14px;
  background: var(--bg-ffffff);
  padding: 14px 18px;
  border-bottom: 1px solid var(--bd-eef2f6);
  border-radius: 16px 16px 0 0;
  width: 100%;
}

.preview-head-left {
  display: flex;
  align-items: center;
  gap: 12px;
  min-width: 0;
}

.preview-head-icon-box {
  width: 38px;
  height: 38px;
  border-radius: 10px;
  background: var(--bg-eff6ff);
  color: var(--tx-2563eb);
  border: 1px solid var(--bd-bfdbfe);
  display: inline-flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.preview-head-text-col {
  display: flex;
  flex-direction: column;
  gap: 2px;
  min-width: 0;
}

.preview-head-title-row {
  display: flex;
  align-items: center;
  gap: 8px;
  flex-wrap: wrap;
}

.preview-head-title {
  font-size: 0.95rem;
  font-weight: 700;
  color: var(--tx-0f172a);
  line-height: 1.3;
}

.preview-head-cat-badge {
  font-size: 0.75rem;
  font-weight: 700;
  color: var(--tx-2563eb);
  background: var(--bg-eff6ff);
  border: 1px solid var(--bd-dbeafe);
  padding: 1px 8px;
  border-radius: 6px;
  white-space: nowrap;
}

.preview-head-subtitle {
  font-size: 0.78rem;
  color: var(--tx-64748b);
  line-height: 1.4;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  max-width: 460px;
}

.preview-head-right {
  display: flex;
  align-items: center;
  gap: 12px;
  flex-wrap: wrap;
  margin-left: auto;
}

.preview-project-info {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 0.88rem;
  color: var(--tx-0f172a);
}

.preview-info-dot {
  width: 10px;
  height: 10px;
  border-radius: 50%;
  flex-shrink: 0;
}

.preview-badge-stalls {
  font-size: 0.76rem;
  font-weight: 700;
  color: var(--tx-1d4ed8);
  background: var(--bg-dbeafe);
  padding: 3px 9px;
  border-radius: 6px;
  border: 1px solid var(--bd-bfdbfe);
}

.preview-legend-row {
  display: flex;
  align-items: center;
  gap: 14px;
  font-size: 0.78rem;
  color: var(--tx-475569);
}

.preview-leg-item {
  display: flex;
  align-items: center;
  gap: 5px;
}

.leg-dot {
  width: 9px;
  height: 9px;
  border-radius: 2px;
  display: inline-block;
}

.leg-dot.is-selected {
  background: #2563eb;
  box-shadow: 0 0 0 2px rgba(37, 99, 235, 0.2);
}

.leg-dot.is-available {
  background: #10b981;
}

.leg-dot.is-reserved {
  background: #f59e0b;
}

/* Selected Stall Checkmark Icon */
.selected-stall-pin-marker {
  pointer-events: none;
  z-index: 25;
}

/* Return to Booked Stalls Pill Button (Top bar) */
.preview-booked-pill-btn {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  background: var(--bg-ffffff);
  border: 1.5px solid var(--bd-bfdbfe);
  border-radius: 10px;
  padding: 5px 12px;
  cursor: pointer;
  transition: all 0.2s ease;
  font-family: inherit;
  box-shadow: 0 1px 3px rgba(37, 99, 235, 0.08);
}

.preview-booked-pill-btn:hover {
  border-color: #2563eb;
  background: var(--bg-eff6ff);
  transform: translateY(-1px);
  box-shadow: 0 4px 10px rgba(37, 99, 235, 0.16);
}

.preview-booked-pill-btn.is-active {
  border-color: var(--bd-93c5fd);
  background: var(--bg-f8fafc);
}

.return-hint-tag {
  background: #2563eb;
  color: #ffffff;
  font-size: 0.75rem;
  font-weight: 700;
  padding: 2px 8px;
  border-radius: 6px;
  display: inline-flex;
  align-items: center;
  margin-left: 4px;
  box-shadow: 0 2px 6px rgba(37, 99, 235, 0.3);
}

/* Floating Return Button over Map */
.map-floating-return-btn {
  position: absolute;
  bottom: 16px;
  left: 50%;
  transform: translateX(-50%);
  z-index: 25;
  background: rgba(15, 23, 42, 0.92);
  color: #ffffff;
  backdrop-filter: blur(10px);
  border: 1.5px solid var(--bd-rgba-255-255-255-0p2);
  padding: 8px 18px;
  border-radius: 30px;
  font-size: 0.84rem;
  display: inline-flex;
  align-items: center;
  gap: 8px;
  cursor: pointer;
  box-shadow: 0 10px 25px rgba(0, 0, 0, 0.25);
  transition: all 0.2s cubic-bezier(0.16, 1, 0.3, 1);
  font-family: inherit;
  white-space: nowrap;
}

.map-floating-return-btn:hover {
  background: #2563eb;
  border-color: #60a5fa;
  transform: translateX(-50%) translateY(-2px);
  box-shadow: 0 12px 30px rgba(37, 99, 235, 0.4);
}

.pulse-indicator-dot {
  width: 9px;
  height: 9px;
  border-radius: 50%;
  background: #60a5fa;
  box-shadow: 0 0 0 0 rgba(96, 165, 250, 0.7);
  animation: pulseBlue 1.6s infinite;
}

@keyframes pulseBlue {
  0% {
    transform: scale(0.95);
    box-shadow: 0 0 0 0 rgba(96, 165, 250, 0.7);
  }
  70% {
    transform: scale(1);
    box-shadow: 0 0 0 8px rgba(96, 165, 250, 0);
  }
  100% {
    transform: scale(0.95);
    box-shadow: 0 0 0 0 rgba(96, 165, 250, 0);
  }
}

.floating-stall-codes-badge {
  background: var(--bg-rgba-255-255-255-0p2);
  border-radius: 4px;
  padding: 1px 6px;
  font-size: 0.76rem;
  font-weight: 700;
}

/* Beacon Pin for Booked Project when viewing other projects */
.booked-project-beacon-marker {
  cursor: pointer;
  z-index: 28;
}

.beacon-pin-bubble {
  position: relative;
  background: #2563eb;
  color: #ffffff;
  padding: 5px 12px;
  border-radius: 20px;
  font-size: 0.76rem;
  font-weight: 700;
  display: inline-flex;
  align-items: center;
  gap: 5px;
  box-shadow: 0 6px 16px rgba(37, 99, 235, 0.45);
  border: 1.5px solid var(--bd-ffffff);
  white-space: nowrap;
  transition: transform 0.2s ease;
}

.beacon-pin-bubble:hover {
  transform: scale(1.08);
  background: #1d4ed8;
}

.beacon-pulse-ring {
  position: absolute;
  top: 50%;
  left: 12px;
  transform: translate(-50%, -50%);
  width: 18px;
  height: 18px;
  border-radius: 50%;
  border: 2px solid #60a5fa;
  animation: beaconPulse 1.8s infinite;
  pointer-events: none;
}

@keyframes beaconPulse {
  0% {
    transform: translate(-50%, -50%) scale(0.6);
    opacity: 1;
  }
  100% {
    transform: translate(-50%, -50%) scale(2.4);
    opacity: 0;
  }
}

@keyframes checkPop {
  0% { transform: scale(0.3); opacity: 0; }
  70% { transform: scale(1.15); opacity: 1; }
  100% { transform: scale(1); opacity: 1; }
}


/* ===== หน้าข้อมูลตลาด: แถบเครื่องมือ + แผงสไลด์จากขวา ===== */
.ov-rail { position: absolute; top: 44%; right: 0; z-index: 120; display: flex; flex-direction: column; gap: 4px; padding: 5px; border: 1px solid var(--bd-rgba-226-232-240-0p95); border-radius: 14px; background: var(--bg-rgba-255-255-255-0p96); backdrop-filter: blur(8px); box-shadow: 0 8px 24px rgba(15, 23, 42, 0.14); transform: translateY(-50%); transition: right 0.26s cubic-bezier(0.16, 1, 0.3, 1), opacity 0.2s ease; }
.ov-rail.is-shifted { right: min(380px, 92vw); }
.ov-rail.is-shifted-dashboard { right: min(380px, 92vw); }
.ov-rail.is-hidden { opacity: 0; pointer-events: none; }
.ov-rail-btn { position: relative; width: 42px; height: 42px; display: grid; place-items: center; border: 0; border-radius: 10px; background: transparent; color: var(--tx-475569); cursor: pointer; transition: background 0.15s ease, color 0.15s ease; }
.ov-rail-btn:hover { background: var(--bg-eff6ff); color: var(--tx-2563eb); }
.ov-rail-btn.active { background: var(--bg-eff6ff); color: var(--tx-2563eb); }
.ov-rail-btn:focus-visible { outline: 2px solid var(--bd-93c5fd); outline-offset: 1px; }
.ov-rail-badge { position: absolute; top: 3px; right: 3px; min-width: 16px; height: 16px; padding: 0 4px; border-radius: 9999px; background: #2563eb; color: #fff; font-size: 0.75rem; font-weight: 700; line-height: 16px; text-align: center; }

.ov-slide { overscroll-behavior: contain; position: absolute; z-index: 140; top: 0; right: 0; bottom: 0; box-sizing: border-box; width: min(380px, 92vw); display: flex; flex-direction: column; background: var(--bg-ffffff); border-left: 1px solid var(--bd-e2e8f0); box-shadow: -12px 0 32px rgba(15, 23, 42, 0.16); }
.ov-slide.is-dashboard-panel { width: min(380px, 92vw); }
.ov-slide.is-dashboard-panel .ov-slide-head { padding: 16px 18px; }
.ov-slide.is-dashboard-panel .ov-slide-title h3 { font-size: 1.125rem; font-weight: 700; color: var(--tx-0f172a); }
.ov-slide.is-dashboard-panel .ov-slide-title p { font-size: 0.875rem; color: var(--tx-64748b); font-weight: 500; margin-top: 2px; }
.ov-slide.is-dashboard-panel .ov-slide-body { padding: 16px 18px; }
.ov-slide-enter-active, .ov-slide-leave-active { transition: transform 0.26s cubic-bezier(0.16, 1, 0.3, 1), opacity 0.2s ease; }
.ov-slide-enter-from, .ov-slide-leave-to { transform: translateX(100%); opacity: 0.6; }
.ov-slide-head { flex: 0 0 auto; display: flex; align-items: center; justify-content: space-between; gap: 12px; padding: 16px 18px; border-bottom: 1px solid var(--bd-f1f5f9); }
.ov-slide-head h3 { margin: 0; color: var(--tx-0f172a); font-size: 1.125rem; font-weight: 700; }
.ov-slide-close { width: 34px; height: 34px; display: grid; place-items: center; border: 1px solid var(--bd-e2e8f0); border-radius: 8px; background: var(--bg-ffffff); color: var(--tx-475569); cursor: pointer; }
.ov-slide-close:hover { background: var(--bg-f1f5f9); }
.ov-slide-body { flex: 1 1 auto; min-height: 0; overflow-y: auto; overflow-x: hidden; padding: 16px 18px; display: flex; flex-direction: column; gap: 16px; -webkit-overflow-scrolling: touch; }
.ov-slide-foot { flex: 0 0 auto; display: flex; justify-content: space-between; gap: 8px; padding: 12px 18px; border-top: 1px solid var(--bd-f1f5f9); }
.ov-search { position: relative; display: flex; align-items: center; }
.ov-search > svg { position: absolute; left: 13px; color: var(--tx-94a3b8); pointer-events: none; }
.ov-search input { box-sizing: border-box; width: 100%; height: 42px; padding: 0 38px 0 38px; border: 1px solid var(--bd-cbd5e1); border-radius: 10px; color: var(--tx-1e293b); font-family: inherit; font-size: 0.9375rem; font-weight: 500; outline: none; }
.ov-search input:focus { border-color: #2563eb; box-shadow: 0 0 0 3px rgba(37, 99, 235, 0.12); }
.ov-search-clear { position: absolute; right: 8px; width: 26px; height: 26px; display: grid; place-items: center; border: 0; border-radius: 6px; background: transparent; color: var(--tx-64748b); cursor: pointer; }
.ov-search-clear:hover { background: var(--bg-f1f5f9); }
.ov-applied-search { display: flex; align-items: center; justify-content: space-between; gap: 10px; margin-top: -8px; padding: 8px 10px; border-radius: 9px; background: var(--bg-eff6ff); color: var(--tx-475569); font-size: 0.8125rem; font-weight: 500; }
.ov-applied-search strong { color: var(--tx-1d4ed8); font-weight: 700; overflow-wrap: anywhere; }
.ov-applied-search button { flex: 0 0 auto; padding: 0; border: 0; background: transparent; color: var(--tx-2563eb); font-family: inherit; font-size: 0.8125rem; font-weight: 600; cursor: pointer; }
.ov-applied-search button:hover { text-decoration: underline; }
.ov-hint { margin: 0; color: var(--tx-64748b); font-size: 0.875rem; font-weight: 500; line-height: 1.5; }
.ov-search-summary { margin: 0; color: var(--tx-64748b); font-size: 0.875rem; font-weight: 500; }
.ov-search-summary strong { color: var(--tx-1e293b); font-weight: 700; }
.ov-result-list, .ov-legend { list-style: none; margin: 0; padding: 0; display: flex; flex-direction: column; gap: 6px; }
.ov-result { width: 100%; display: flex; align-items: center; gap: 12px; padding: 10px 12px; border: 1px solid var(--bd-eef2f6); border-radius: 12px; background: var(--bg-ffffff); text-align: left; cursor: pointer; font-family: inherit; }
.ov-result:hover { background: var(--bg-f8fafc); border-color: var(--bd-dbeafe); }
.ov-result.focused { border-color: var(--bd-bfdbfe); background: var(--bg-f8fbff); }
.ov-result.active { border-color: var(--bd-93c5fd); background: var(--bg-eff6ff); box-shadow: 0 0 0 2px rgba(37, 99, 235, 0.08); }
.ov-result-text, .ov-legend-text { display: flex; flex-direction: column; min-width: 0; }
.ov-result-text strong, .ov-legend-text strong { color: var(--tx-1e293b); font-size: 0.9375rem; font-weight: 600; }
.ov-result-text small, .ov-legend-text small { color: var(--tx-64748b); font-size: 0.8125rem; font-weight: 500; }
.ov-result-text .ov-result-reference { color: var(--tx-2563eb); font-weight: 600; overflow-wrap: anywhere; }
.ov-result-text mark { padding: 0 1px; border-radius: 3px; background: var(--bg-fef3c7); color: inherit; }
.ov-result-badges { display: flex; flex-wrap: wrap; gap: 4px; margin: 3px 0 2px; }
.ov-result-badge { display: inline-flex; align-items: center; min-height: 20px; padding: 2px 7px; border-radius: 6px; background: var(--bg-f1f5f9); color: var(--tx-475569); font-size: 0.75rem; font-weight: 600; line-height: 1.2; }
.ov-result-badge.is-booking { background: var(--bg-dbeafe); color: var(--tx-1d4ed8); }
.ov-result-badge.is-application { background: var(--bg-f3e8ff); color: var(--tx-7e22ce); }
.ov-result-badge.is-stall { background: var(--bg-dcfce7); color: var(--tx-15803d); }
.ov-result-badge.is-project { background: var(--bg-f1f5f9); color: var(--tx-475569); }
.ov-swatch { flex: 0 0 auto; width: 16px; height: 16px; border-radius: 5px; box-shadow: inset 0 0 0 1px rgba(15, 23, 42, 0.12); }
.ov-legend-item { width: 100%; display: flex; align-items: center; gap: 12px; padding: 8px 10px; border: 1px solid transparent; border-radius: 10px; background: var(--bg-f8fafc); text-align: left; font-family: inherit; cursor: pointer; transition: background 0.15s ease, border-color 0.15s ease, opacity 0.15s ease; }
.ov-legend-item:hover { background: var(--bg-f1f5f9); }
.ov-legend-item:focus-visible { outline: 2px solid var(--bd-93c5fd); outline-offset: 1px; }
.ov-legend-item.active { border-color: var(--bd-93c5fd); background: var(--bg-eff6ff); }
.ov-legend-item.dimmed { opacity: 0.55; }
.ov-legend-check { flex: 0 0 auto; margin-left: auto; color: var(--tx-2563eb); }
.ov-legend-clear { margin: 0 0 10px; padding: 0; border: 0; background: none; color: var(--tx-2563eb); font-family: inherit; font-size: 0.875rem; font-weight: 600; cursor: pointer; }
.ov-legend-clear:hover { text-decoration: underline; }
.ov-field { display: flex; flex-direction: column; gap: 6px; }
.ov-label { color: var(--tx-475569); font-size: 0.875rem; font-weight: 600; }
.ov-result-summary { margin: 0; padding: 10px 12px; border-radius: 10px; background: var(--bg-eff6ff); color: var(--tx-475569); font-size: 0.875rem; font-weight: 500; }
.ov-result-summary strong { color: var(--tx-2563eb); font-weight: 700; }
.ov-section h4 { margin: 0 0 10px; color: var(--tx-1e293b); font-size: 1rem; font-weight: 700; }
.ov-section .ov-hint { margin-bottom: 10px; }
.ov-basemaps { display: grid; grid-template-columns: repeat(3, minmax(0, 1fr)); gap: 10px; }
.ov-basemap { display: flex; flex-direction: column; align-items: center; gap: 6px; padding: 6px 4px 8px; border: 2px solid transparent; border-radius: 12px; background: var(--bg-f8fafc); color: var(--tx-475569); font-family: inherit; font-size: 0.8125rem; font-weight: 600; cursor: pointer; }
.ov-basemap:hover { background: var(--bg-f1f5f9); }
.ov-basemap.active { border-color: #2563eb; background: var(--bg-eff6ff); color: var(--tx-2563eb); }
.ov-basemap-thumb { width: 100%; aspect-ratio: 1; border-radius: 8px; background-size: cover; }
.ov-basemap-thumb.is-standard { background: linear-gradient(135deg, transparent 46%, #fde68a 47% 53%, transparent 54%), linear-gradient(45deg, transparent 46%, var(--bg-ffffff) 47% 53%, transparent 54%), var(--bg-e8eef6); }
.ov-basemap-thumb.is-satellite { background: radial-gradient(circle at 30% 35%, #3f6b3a 0 22%, transparent 23%), radial-gradient(circle at 70% 65%, #6b5a3a 0 24%, transparent 25%), #1f3a2a; }
.ov-basemap-thumb.is-dark { background: linear-gradient(135deg, transparent 46%, #475569 47% 53%, transparent 54%), linear-gradient(45deg, transparent 46%, #334155 47% 53%, transparent 54%), #0f172a; }
@media (max-width: 720px) {
  .ov-rail.is-shifted { opacity: 0; pointer-events: none; }
}

/* ===== ปรับโฉมแผงสไลด์ / แถบเครื่องมือ ===== */
.ov-rail { padding: 6px; gap: 4px; border-right: 0; border-radius: 14px 0 0 14px; }
.ov-rail-btn { box-sizing: border-box; flex: 0 0 40px; width: 40px; height: 40px; border-radius: 8px; }
.ov-rail-btn { display: flex; align-items: center; justify-content: center; }
.ov-rail-label { font-size: 0.75rem; font-weight: 600; line-height: 1; letter-spacing: 0.01em; }
.ov-rail-badge { top: 2px; right: 2px; }
.ov-slide { border-radius: 18px 0 0 18px; }
.ov-slide-head { gap: 12px; padding: 18px 18px 14px; background: linear-gradient(180deg, var(--bg-f8fbff) 0%, var(--bg-ffffff) 100%); }
.ov-slide-icon { flex: 0 0 auto; width: 40px; height: 40px; display: grid; place-items: center; border-radius: 12px; background: var(--bg-eff6ff); color: var(--tx-2563eb); }
.ov-slide-title { flex: 1 1 auto; min-width: 0; }
.ov-slide-title h3 { margin: 0; }
.ov-slide-title p { margin: 2px 0 0; color: var(--tx-64748b); font-size: 0.8125rem; font-weight: 500; }
.ov-slide-body { gap: 18px; padding: 18px; }
.ov-slide-foot { padding: 14px 18px; background: var(--bg-ffffff); box-shadow: 0 -6px 16px rgba(15, 23, 42, 0.04); }
.ov-slide-foot > * { flex: 1 1 0; }
.ov-search input { height: 46px; border-radius: 12px; background: var(--bg-f8fafc); }
.ov-search input:focus { background: var(--bg-ffffff); }
.ov-empty { display: flex; flex-direction: column; align-items: center; gap: 6px; padding: 32px 16px; border: 1px dashed var(--bd-dbe3ee); border-radius: 14px; background: var(--bg-f8fafc); text-align: center; }
.ov-empty-icon { width: 48px; height: 48px; display: grid; place-items: center; border-radius: 50%; background: var(--bg-eff6ff); color: var(--tx-2563eb); }
.ov-empty strong { color: var(--tx-1e293b); font-size: 0.9375rem; font-weight: 700; }
.ov-empty p { margin: 0; color: var(--tx-64748b); font-size: 0.875rem; font-weight: 500; line-height: 1.5; }
.ov-result { border-radius: 14px; transition: background 0.15s ease, border-color 0.15s ease, transform 0.15s ease; }
.ov-result:hover { transform: translateX(-2px); }
.ov-chevron { flex: 0 0 auto; margin-left: auto; color: var(--tx-94a3b8); }
.ov-section { padding: 14px; border: 1px solid var(--bd-eef2f6); border-radius: 16px; background: var(--bg-ffffff); box-shadow: 0 1px 3px rgba(15, 23, 42, 0.03); }
.ov-basemaps { gap: 8px; }
.ov-basemap { gap: 8px; padding: 6px 6px 10px; border-radius: 14px; transition: transform 0.15s ease, background 0.15s ease, border-color 0.15s ease; }
.ov-basemap:hover { transform: translateY(-1px); }
.ov-basemap-thumb { position: relative; border-radius: 10px; box-shadow: inset 0 0 0 1px rgba(15, 23, 42, 0.08); }
.ov-basemap-check { position: absolute; top: 6px; right: 6px; width: 22px; height: 22px; padding: 4px; box-sizing: border-box; border-radius: 50%; background: #2563eb; color: #fff; box-shadow: 0 2px 6px rgba(37, 99, 235, 0.4); }
.ov-legend { gap: 6px; }
.ov-legend-item { padding: 10px 12px; border-radius: 12px; }
.ov-swatch { width: 20px; height: 20px; border-radius: 7px; }
.ov-legend-count { flex: 0 0 auto; margin-left: auto; min-width: 26px; padding: 2px 8px; border-radius: 9999px; background: var(--bg-e2e8f0); color: var(--tx-475569); font-size: 0.8125rem; font-weight: 700; text-align: center; }
.ov-field { padding: 12px; border: 1px solid var(--bd-eef2f6); border-radius: 14px; background: var(--bg-f8fafc); }
.ov-result-summary { padding: 12px 14px; border-radius: 12px; }
@media (max-width: 768px) {
  .ov-slide,
  .ov-slide.is-dashboard-panel {
    width: 100% !important;
    width: 100vw !important;
    max-width: 100% !important;
    left: 0 !important;
    right: 0 !important;
    border-left: 0 !important;
    border-radius: 0 !important;
    box-shadow: none !important;
  }
  .ov-rail.is-shifted,
  .ov-rail.is-shifted-dashboard {
    display: none !important;
    opacity: 0 !important;
    pointer-events: none !important;
  }
}

/* เปิดแผงสไลด์: ดันปุ่มซูม/เข็มทิศ/กึ่งกลาง และเครดิตแผนที่ ออกมาทางซ้ายของแผง ไม่ให้ถูกบัง */
.maplibre-overview-wrapper .maplibregl-ctrl-bottom-right { transition: transform 0.26s cubic-bezier(0.16, 1, 0.3, 1); }
.maplibre-overview-wrapper.has-ov-panel .maplibregl-ctrl-bottom-right { transform: translateX(calc(-1 * min(380px, 92vw))); }
.maplibre-overview-wrapper.has-ov-panel .map-floating-top-controls { transform: translateX(calc(-1 * min(380px, 92vw))); }
.maplibre-overview-wrapper.has-ov-dashboard .maplibregl-ctrl-bottom-right { transform: translateX(calc(-1 * min(380px, 92vw))); }
.maplibre-overview-wrapper.has-ov-dashboard .map-floating-top-controls { transform: translateX(calc(-1 * min(380px, 92vw))); }
@media (max-width: 768px) {
  .maplibre-overview-wrapper.has-ov-panel .maplibregl-ctrl-bottom-right,
  .maplibre-overview-wrapper.has-ov-panel .map-floating-top-controls,
  .maplibre-overview-wrapper.has-ov-dashboard .maplibregl-ctrl-bottom-right,
  .maplibre-overview-wrapper.has-ov-dashboard .map-floating-top-controls { transform: none; opacity: 0; pointer-events: none; }
}

/* ===== ชั้นข้อมูล: รายการสีโครงการแบบเรียบ อ่านง่าย ===== */
.ov-section { padding: 0; border: 0; border-radius: 0; background: transparent; box-shadow: none; }
.ov-section + .ov-section { padding-top: 18px; border-top: 1px solid var(--bd-f1f5f9); }
.ov-section h4 { margin: 0 0 4px; font-size: 1rem; }
.ov-section .ov-hint { margin: 0 0 12px; font-size: 0.8125rem; line-height: 1.5; color: var(--tx-94a3b8); }
.ov-legend { gap: 2px; margin: 0 -8px; }
.ov-legend-item { padding: 10px 8px; gap: 12px; align-items: center; border: 0; border-radius: 10px; background: transparent; }
.ov-legend-item:hover { background: var(--bg-f8fafc); }
.ov-legend-item.active { background: var(--bg-eff6ff); }
.ov-legend-item.dimmed { opacity: 0.5; }
.ov-legend .ov-swatch { width: 14px; height: 14px; border-radius: 50%; box-shadow: 0 0 0 3px rgba(148, 163, 184, 0.18); }
.ov-legend-text { flex: 1 1 auto; min-width: 0; gap: 1px; }
.ov-legend-text strong { color: var(--tx-1e293b); font-size: 0.9375rem; font-weight: 600; line-height: 1.35; }
.ov-legend-text small { color: var(--tx-94a3b8); font-size: 0.8125rem; font-weight: 500; white-space: nowrap; overflow: hidden; text-overflow: ellipsis; }
.ov-legend-count { min-width: 0; padding: 0; background: none; color: var(--tx-94a3b8); font-size: 0.875rem; font-weight: 600; }
.ov-legend-item.active .ov-legend-count { display: none; }
.ov-legend-clear { margin: 0 0 8px; }

/* ตัวกรอง: ไม่ใส่กรอบครอบแต่ละช่อง */
.ov-field { padding: 0; border: 0; border-radius: 0; background: transparent; }
.ov-filter-section { display: flex; flex-direction: column; gap: 10px; }
.ov-kpi-filter-chip { display: flex; align-items: center; justify-content: space-between; gap: 10px; padding: 9px 10px 9px 12px; border: 1px solid var(--bd-bfdbfe); border-radius: 10px; background: var(--bg-eff6ff); color: var(--tx-475569); font-size: 0.8125rem; font-weight: 500; }
.ov-kpi-filter-chip strong { color: var(--tx-1d4ed8); font-weight: 700; }
.ov-kpi-filter-chip button { flex: 0 0 auto; width: 28px; height: 28px; display: grid; place-items: center; padding: 0; border: 0; border-radius: 7px; background: transparent; color: var(--tx-64748b); cursor: pointer; }
.ov-kpi-filter-chip button:hover { background: var(--bg-dbeafe); color: var(--tx-1d4ed8); }
.ov-kpi-filter-chip button:focus-visible { outline: 2px solid var(--bd-93c5fd); outline-offset: 2px; }
.ov-filter-section + .ov-filter-section { padding-top: 16px; border-top: 1px solid var(--bd-f1f5f9); }
.ov-filter-section h4 { margin: 0; color: var(--tx-1e293b); font-size: 0.9375rem; font-weight: 700; }
.ov-filter-grid { display: grid; grid-template-columns: 1fr; gap: 12px; }
.ov-filter-grid .ov-label { min-height: 20px; line-height: 1.4; }
@media (max-width: 420px) {
  .ov-filter-grid { grid-template-columns: 1fr; }
  .ov-filter-grid .ov-label { min-height: 0; }
}

.ov-tile { flex: 0 0 auto; width: 36px; height: 36px; display: grid; place-items: center; border-radius: 10px; color: #fff; font-size: 0.9375rem; font-weight: 700; text-shadow: 0 1px 1px rgba(0, 0, 0, 0.18); }
.ov-legend-item.active .ov-tile { box-shadow: 0 0 0 2px #fff, 0 0 0 4px #2563eb; }

.ov-section-head { display: flex; align-items: center; justify-content: space-between; gap: 10px; margin-bottom: 4px; }
.ov-section-head h4 { margin: 0; }
.ov-layer-toggle-list { overflow: visible; border: 0; border-radius: 11px; background: transparent; }
.ov-layer-toggle { width: 100%; display: flex; align-items: center; gap: 8px; padding: 9px 10px; border: 1px solid var(--bd-e2e8f0); border-radius: 11px; background: var(--bg-ffffff); color: var(--tx-64748b); font-family: inherit; text-align: left; cursor: pointer; transition: border-color 0.15s ease, background 0.15s ease, color 0.15s ease; }
.ov-layer-toggle:hover { background: var(--bg-f8fafc); }
.ov-layer-toggle:focus-visible { outline: 2px solid var(--bd-93c5fd); outline-offset: 2px; }
.ov-layer-toggle.active { border-color: var(--bd-bfdbfe); color: var(--tx-1e3a8a); background: var(--bg-f8fbff); }
.ov-layer-icon { flex: 0 0 auto; width: 28px; height: 28px; display: grid; place-items: center; border-radius: 8px; background: var(--bg-f1f5f9); color: var(--tx-64748b); }
.ov-layer-toggle.active .ov-layer-icon { background: var(--bg-dbeafe); color: var(--tx-2563eb); }
.ov-layer-text { flex: 1 1 auto; min-width: 0; display: flex; flex-direction: column; gap: 1px; }
.ov-layer-text strong { color: var(--tx-1e293b); font-size: 0.9375rem; font-weight: 600; }
.ov-layer-text small { color: var(--tx-64748b); font-size: 0.8125rem; font-weight: 500; }
.ov-layer-switch { position: relative; flex: 0 0 auto; width: 30px; height: 18px; border-radius: 9999px; background: var(--bg-cbd5e1); transition: background 0.18s ease; }
.ov-layer-switch span { position: absolute; top: 3px; left: 3px; width: 12px; height: 12px; border-radius: 50%; background: var(--bg-ffffff); box-shadow: 0 1px 2px rgba(15, 23, 42, 0.25); transition: transform 0.18s ease; }
.ov-layer-toggle.active .ov-layer-switch { background: #2563eb; }
.ov-layer-toggle.active .ov-layer-switch span { transform: translateX(12px); }
.ov-chip-clear { display: inline-flex; align-items: center; gap: 5px; height: 30px; padding: 0 12px; border: 1px solid var(--bd-fca5a5); border-radius: 9999px; background: var(--bg-ffffff); color: var(--tx-dc2626); font-family: inherit; font-size: 0.8125rem; font-weight: 600; cursor: pointer; transition: background 0.15s ease; }
.ov-chip-clear:hover { background: var(--bg-fef2f2); }
.ov-chip-clear:focus-visible { outline: 2px solid var(--bd-fca5a5); outline-offset: 2px; }

/* ===== สีโครงการ: เลือกได้หลายหมวด ===== */
.ov-legend { gap: 0; margin: 0; border: 1px solid var(--bd-eef2f6); border-radius: 14px; overflow: hidden; }
.ov-legend li + li { border-top: 1px solid var(--bd-f1f5f9); }
.ov-legend-item { padding: 12px 14px; border-radius: 0; background: var(--bg-ffffff); }
.ov-legend-item:hover { background: var(--bg-f8fafc); }
.ov-legend-item.active { background: var(--bg-f5f9ff); }
.ov-legend-item.dimmed { opacity: 0.55; }
.ov-legend-item.active .ov-tile { box-shadow: none; }
.ov-check { flex: 0 0 auto; width: 20px; height: 20px; display: grid; place-items: center; border: 2px solid var(--bd-cbd5e1); border-radius: 6px; color: #fff; transition: background 0.15s ease, border-color 0.15s ease; }
.ov-legend-item.active .ov-check { background: #2563eb; border-color: #2563eb; }
.ov-tile { width: 34px; height: 34px; }

/* เพิ่มระยะขอบซ้าย-ขวาของแผงสไลด์ */
.ov-slide-head { padding-left: 24px; padding-right: 24px; }
.ov-slide-body { padding-left: 24px; padding-right: 24px; }
.ov-slide-foot { padding-left: 24px; padding-right: 24px; }
@media (max-width: 768px) {
  .ov-slide-head, .ov-slide-body, .ov-slide-foot { padding-left: 16px; padding-right: 16px; }
}

/* ซ่อนปุ่มล้างของเบราว์เซอร์ในช่อง type=search เพราะมีปุ่มล้างของระบบอยู่แล้ว */
.ov-search input::-webkit-search-cancel-button,
.ov-search input::-webkit-search-decoration { -webkit-appearance: none; appearance: none; display: none; }

/* ปุ่มล้างตัวกรองในชั้นข้อมูล: ย่อให้เล็กลง */
.ov-chip-clear { height: 26px; padding: 0 10px; gap: 4px; font-size: 0.8125rem; font-weight: 600; }
.ov-chip-clear svg { width: 11px; height: 11px; }

/* การ์ดโครงการบนแผนที่ (ตอนชี้): ขยายตัวอักษรให้อ่านง่ายตามมาตรฐานระบบ */
.mlp-body { padding: 16px 18px 14px; }
.mlp-title strong { font-size: 1.1875rem; }
.mlp-pill { font-size: 0.875rem; padding: 4px 12px; }
.mlp-building-badge { font-size: 0.8125rem; }
.mlp-sub { font-size: 0.9375rem; line-height: 1.45; }
.mlp-tile { padding: 10px 12px; }
.mlp-tile-label { font-size: 0.875rem; }
.mlp-tile-value { font-size: 1.3125rem; }
.mlp-tile-value small { font-size: 0.875rem; }
.mlp-tile-note { font-size: 0.8125rem; }
.mlp-occ-head { font-size: 0.875rem; }
.mlp-occ-legend { font-size: 0.9375rem; gap: 14px; }
.mlp-occ-legend i { width: 10px; height: 10px; }
.mlp-hint { font-size: 0.875rem; line-height: 1.5; }
.mlp-queue-title { font-size: 0.9375rem; }
.mlp-queue-desc { font-size: 0.8125rem; }

/* การ์ดโครงการบนแผนที่: เพิ่มระยะห่างแนวตั้งให้โปร่งขึ้น */
.mlp-body { padding: 20px 20px 18px; }
.mlp-head { margin-bottom: 4px; }
.mlp-sub { margin-top: 6px; }
.mlp-tiles { margin-top: 16px; gap: 12px; }
.mlp-tile { padding: 14px 14px 12px; }
.mlp-occ { margin-top: 20px; }
.mlp-occ-head { margin-bottom: 8px; }
.mlp-occ-bar { height: 8px; margin: 10px 0 12px; }
.mlp-occ-legend { margin-top: 12px; }
.mlp-hint { margin-top: 18px; padding-top: 14px; }

/* ปุ่ม "ปรับกึ่งกลาง" แบบลอยเดิม ถูกรวมเข้าชุดปุ่มซูมแล้ว ทุกหน้า */
.map-floating-recenter-btn { display: none !important; }

/* แผนที่ฐานเหลือ 2 แบบ: ให้เต็มความกว้าง */
.ov-basemaps { grid-template-columns: repeat(2, minmax(0, 1fr)); }

/* ชั้นข้อมูล: ลดความสูงของตัวเลือกให้เห็นข้อมูลได้มากขึ้น */
.ov-section + .ov-section { padding-top: 14px; }
.ov-basemaps { gap: 10px; }
.ov-basemap { min-width: 0; gap: 7px; padding: 6px 6px 9px; border-width: 1.5px; }
.ov-basemap-thumb { height: 88px; aspect-ratio: auto; overflow: hidden; }
.ov-basemap-thumb::before,
.ov-basemap-thumb::after { content: ''; position: absolute; pointer-events: none; }
.ov-basemap-thumb.is-standard::before { inset: 14% -12%; border-top: 5px solid var(--bd-rgba-255-255-255-0p9); border-bottom: 2px solid rgba(148, 163, 184, 0.28); transform: rotate(28deg); }
.ov-basemap-thumb.is-standard::after { width: 9px; height: 9px; left: 22%; top: 24%; border: 3px solid #2563eb; border-radius: 50% 50% 50% 0; background: var(--bg-ffffff); transform: rotate(-45deg); box-shadow: 42px 38px 0 -2px #fff, 42px 38px 0 1px #64748b; }
.ov-basemap-thumb.is-satellite { background: radial-gradient(circle at 27% 34%, rgba(74, 124, 64, 0.95) 0 17%, transparent 18%), radial-gradient(circle at 70% 64%, rgba(139, 109, 68, 0.85) 0 19%, transparent 20%), linear-gradient(135deg, transparent 43%, var(--bg-rgba-207-213-188-0p42) 44% 49%, transparent 50%), #183c2b; }
.ov-basemap-thumb.is-satellite::before { inset: -12% 48% -12% 42%; background: rgba(160, 181, 151, 0.28); transform: rotate(18deg); }
.ov-basemap-thumb.is-satellite::after { width: 10px; height: 10px; right: 18%; top: 20%; border: 2px solid var(--bd-ffffff); border-radius: 50%; background: #2563eb; box-shadow: 0 2px 6px rgba(15, 23, 42, 0.35); }
.ov-basemap.active { box-shadow: 0 7px 18px rgba(37, 99, 235, 0.13); }
.ov-layer-toggle { position: relative; }
.ov-layer-toggle.active { box-shadow: none; background: var(--bg-f8fbff); }
.ov-label-layer:focus-visible { outline: 2px solid var(--bd-93c5fd); outline-offset: 2px; }
.ov-label-layer-list { display: grid; grid-template-columns: repeat(2, minmax(0, 1fr)); gap: 8px; }
.ov-label-layer { min-width: 0; display: grid; grid-template-columns: auto minmax(0, 1fr) auto; align-items: center; gap: 8px; padding: 9px 10px; border: 1px solid var(--bd-e2e8f0); border-radius: 11px; background: var(--bg-ffffff); color: var(--tx-64748b); font-family: inherit; font-size: 0.8125rem; font-weight: 600; text-align: left; cursor: pointer; transition: border-color 0.15s ease, background 0.15s ease, color 0.15s ease; }
.ov-label-layer:hover { background: var(--bg-f8fafc); }
.ov-label-layer.active { border-color: var(--bd-bfdbfe); background: var(--bg-f8fbff); color: var(--tx-1e3a8a); }
.ov-label-layer-icon { width: 28px; height: 28px; display: grid; place-items: center; border-radius: 8px; background: var(--bg-f1f5f9); color: var(--tx-64748b); }
.ov-label-layer.active .ov-label-layer-icon { background: var(--bg-dbeafe); color: var(--tx-2563eb); }
.ov-mini-switch { position: relative; width: 30px; height: 18px; border-radius: 9999px; background: var(--bg-cbd5e1); transition: background 0.15s ease; }
.ov-mini-switch > span { position: absolute; top: 3px; left: 3px; width: 12px; height: 12px; border-radius: 50%; background: var(--bg-ffffff); box-shadow: 0 1px 2px rgba(15, 23, 42, 0.25); transition: transform 0.15s ease; }
.ov-label-layer.active .ov-mini-switch { background: #2563eb; }
.ov-label-layer.active .ov-mini-switch > span { transform: translateX(12px); }
.ov-legend-actions { display: flex; align-items: center; gap: 8px; }
.ov-selected-count { color: var(--tx-64748b); font-size: 0.75rem; font-weight: 600; }
.ov-show-all { padding: 0; border: 0; background: transparent; color: var(--tx-2563eb); font-family: inherit; font-size: 0.8125rem; font-weight: 600; cursor: pointer; }
.ov-show-all:hover { text-decoration: underline; }
.ov-show-all:focus-visible { outline: 2px solid var(--bd-93c5fd); outline-offset: 3px; border-radius: 3px; }
.ov-legend-item { position: relative; box-shadow: inset 0 0 var(--ov-category-color); transition: background 0.15s ease, box-shadow 0.15s ease, opacity 0.15s ease; }
.ov-legend-item.active { background: color-mix(in srgb, var(--ov-category-color) 9%, var(--bg-ffffff)); box-shadow: inset 4px 0 var(--ov-category-color); }
.ov-legend-item.active .ov-check { border-color: var(--ov-category-color); background: var(--ov-category-color); box-shadow: 0 0 0 3px color-mix(in srgb, var(--ov-category-color) 16%, transparent); }
.ov-legend-item.active .ov-tile { transform: scale(1.03); }

/* หมุดแผงที่คุณจอง: วงกลมกระพริบอยู่กึ่งกลางไอคอนพอดี ไม่ล้นออกนอกแคปซูล */
.beacon-pin-bubble { overflow: hidden; }
.beacon-pin-icon { position: relative; display: inline-grid; place-items: center; flex: 0 0 18px; width: 18px; height: 18px; }
.beacon-pin-icon svg { position: relative; z-index: 1; }
.beacon-pin-icon .beacon-pulse-ring {
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  transform: none;
  border: 2px solid var(--bd-rgba-255-255-255-0p85);
  animation: beaconPulseCentered 1.8s ease-out infinite;
}
@keyframes beaconPulseCentered {
  0% { transform: scale(0.7); opacity: 0.95; }
  100% { transform: scale(1.5); opacity: 0; }
}
</style>
