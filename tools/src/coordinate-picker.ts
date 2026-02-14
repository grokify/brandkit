/**
 * Coordinate Picker - TypeScript Library
 * A tool for picking coordinates from images to create SVG polygons.
 */

// ============================================================================
// Types
// ============================================================================

export interface Point {
  x: number;
  y: number;
}

export interface Pin extends Point {
  id: number;
  displayX: number;
  displayY: number;
}

export interface Shape {
  color: string;
  pinIds: number[];
}

export interface CoordinatePickerState {
  pins: Pin[];
  shapes: Shape[];
  activeShapeIndex: number;
  pinIdCounter: number;
  imageWidth: number;
  imageHeight: number;
  zoomLevel: number;
}

export interface CoordinatePickerConfig {
  shapeColors: string[];
  zoomStep: number;
  minZoom: number;
  maxZoom: number;
  pinProximityThreshold: number;
}

// ============================================================================
// Default Configuration
// ============================================================================

export const DEFAULT_CONFIG: CoordinatePickerConfig = {
  shapeColors: [
    '#e74c3c', // red
    '#3498db', // blue
    '#2ecc71', // green
    '#f39c12', // orange
    '#9b59b6', // purple
    '#1abc9c', // teal
    '#e91e63', // pink
    '#00bcd4', // cyan
    '#ff5722', // deep orange
    '#8bc34a', // light green
  ],
  zoomStep: 0.25,
  minZoom: 0.5,
  maxZoom: 4,
  pinProximityThreshold: 5,
};

// ============================================================================
// State Management
// ============================================================================

export function createInitialState(numShapes: number = 4, config: CoordinatePickerConfig = DEFAULT_CONFIG): CoordinatePickerState {
  const shapes: Shape[] = [];
  for (let i = 0; i < numShapes; i++) {
    shapes.push({
      color: config.shapeColors[i % config.shapeColors.length],
      pinIds: [],
    });
  }

  return {
    pins: [],
    shapes,
    activeShapeIndex: 0,
    pinIdCounter: 0,
    imageWidth: 400,
    imageHeight: 400,
    zoomLevel: 1,
  };
}

// ============================================================================
// Pin Management
// ============================================================================

export function findNearbyPin(
  pins: Pin[],
  x: number,
  y: number,
  threshold: number = DEFAULT_CONFIG.pinProximityThreshold
): Pin | undefined {
  return pins.find(p => Math.abs(p.x - x) < threshold && Math.abs(p.y - y) < threshold);
}

export function createPin(
  state: CoordinatePickerState,
  x: number,
  y: number,
  displayX: number,
  displayY: number
): Pin {
  const pin: Pin = {
    id: state.pinIdCounter++,
    x,
    y,
    displayX,
    displayY,
  };
  state.pins.push(pin);
  return pin;
}

export function deletePin(state: CoordinatePickerState, pinId: number): void {
  // Remove from all shapes
  state.shapes.forEach(shape => {
    const idx = shape.pinIds.indexOf(pinId);
    if (idx !== -1) {
      shape.pinIds.splice(idx, 1);
    }
  });

  // Remove pin from global array
  const pinIdx = state.pins.findIndex(p => p.id === pinId);
  if (pinIdx !== -1) {
    state.pins.splice(pinIdx, 1);
  }
}

export function getPinById(state: CoordinatePickerState, pinId: number): Pin | undefined {
  return state.pins.find(p => p.id === pinId);
}

export function getShapesForPin(state: CoordinatePickerState, pinId: number): Shape[] {
  return state.shapes.filter(s => s.pinIds.includes(pinId));
}

export function isPinUsedByMultipleShapes(state: CoordinatePickerState, pinId: number): boolean {
  let count = 0;
  for (const shape of state.shapes) {
    if (shape.pinIds.includes(pinId)) {
      count++;
      if (count > 1) return true;
    }
  }
  return false;
}

// ============================================================================
// Shape Management
// ============================================================================

export function addPinToShape(state: CoordinatePickerState, pinId: number, shapeIndex: number): boolean {
  if (shapeIndex < 0 || shapeIndex >= state.shapes.length) return false;
  if (state.shapes[shapeIndex].pinIds.includes(pinId)) return false;

  state.shapes[shapeIndex].pinIds.push(pinId);
  return true;
}

export function removePinFromShape(state: CoordinatePickerState, pinId: number, shapeIndex: number): boolean {
  if (shapeIndex < 0 || shapeIndex >= state.shapes.length) return false;

  const shape = state.shapes[shapeIndex];
  const idx = shape.pinIds.indexOf(pinId);
  if (idx === -1) return false;

  shape.pinIds.splice(idx, 1);
  return true;
}

export function getShapePins(state: CoordinatePickerState, shapeIndex: number): Pin[] {
  if (shapeIndex < 0 || shapeIndex >= state.shapes.length) return [];

  return state.shapes[shapeIndex].pinIds
    .map(id => getPinById(state, id))
    .filter((p): p is Pin => p !== undefined);
}

export function clearShape(state: CoordinatePickerState, shapeIndex: number): number[] {
  if (shapeIndex < 0 || shapeIndex >= state.shapes.length) return [];

  const shape = state.shapes[shapeIndex];
  const removedPinIds = [...shape.pinIds];
  shape.pinIds = [];

  // Return pins that are no longer used by any shape
  return removedPinIds.filter(pinId => {
    const stillUsed = state.shapes.some(s => s.pinIds.includes(pinId));
    if (!stillUsed) {
      const pinIdx = state.pins.findIndex(p => p.id === pinId);
      if (pinIdx !== -1) {
        state.pins.splice(pinIdx, 1);
      }
    }
    return !stillUsed;
  });
}

export function undoLastPinFromShape(state: CoordinatePickerState, shapeIndex: number): { pinId: number; deleted: boolean } | null {
  if (shapeIndex < 0 || shapeIndex >= state.shapes.length) return null;

  const shape = state.shapes[shapeIndex];
  if (shape.pinIds.length === 0) return null;

  const pinId = shape.pinIds.pop()!;
  const stillUsed = state.shapes.some(s => s.pinIds.includes(pinId));

  if (!stillUsed) {
    const pinIdx = state.pins.findIndex(p => p.id === pinId);
    if (pinIdx !== -1) {
      state.pins.splice(pinIdx, 1);
    }
    return { pinId, deleted: true };
  }

  return { pinId, deleted: false };
}

export function setNumShapes(state: CoordinatePickerState, numShapes: number, config: CoordinatePickerConfig = DEFAULT_CONFIG): void {
  // Preserve existing shapes if reducing number
  while (state.shapes.length < numShapes) {
    state.shapes.push({
      color: config.shapeColors[state.shapes.length % config.shapeColors.length],
      pinIds: [],
    });
  }

  // Truncate if reducing (pins from removed shapes become orphaned)
  if (state.shapes.length > numShapes) {
    state.shapes.length = numShapes;
  }

  // Adjust active shape index if needed
  if (state.activeShapeIndex >= numShapes) {
    state.activeShapeIndex = 0;
  }
}

// ============================================================================
// Coordinate Calculations
// ============================================================================

export function calculateImageCoordinates(
  clickX: number,
  clickY: number,
  zoomLevel: number
): Point {
  return {
    x: Math.round(clickX / zoomLevel),
    y: Math.round(clickY / zoomLevel),
  };
}

export function calculateDisplayPosition(
  imageX: number,
  imageY: number,
  zoomLevel: number
): Point {
  return {
    x: imageX * zoomLevel,
    y: imageY * zoomLevel,
  };
}

// ============================================================================
// SVG Generation
// ============================================================================

export interface PolygonOptions {
  fill?: string;
  stroke?: string;
  strokeWidth?: number;
}

export function generatePolygonPoints(pins: Pin[]): string {
  return pins.map(p => `${p.x},${p.y}`).join(' ');
}

export function generatePolygonElement(pins: Pin[], options: PolygonOptions = {}): string {
  if (pins.length < 3) {
    return '';
  }

  const { fill = '#00FF00', stroke, strokeWidth } = options;
  const points = generatePolygonPoints(pins);

  let attrs = `points="${points}" fill="${fill}"`;
  if (stroke) attrs += ` stroke="${stroke}"`;
  if (strokeWidth) attrs += ` stroke-width="${strokeWidth}"`;

  return `<polygon ${attrs}/>`;
}

export function generateShapePolygon(state: CoordinatePickerState, shapeIndex: number, options: PolygonOptions = {}): string {
  const pins = getShapePins(state, shapeIndex);
  return generatePolygonElement(pins, options);
}

export function generateFullSVG(
  state: CoordinatePickerState,
  options: {
    backgroundColor?: string;
    polygonFill?: string;
    includeBackground?: boolean;
  } = {}
): string {
  const {
    backgroundColor = '#000000',
    polygonFill = '#00FF00',
    includeBackground = true,
  } = options;

  const polygons = state.shapes
    .map(shape => {
      const pins = shape.pinIds
        .map(id => getPinById(state, id))
        .filter((p): p is Pin => p !== undefined);

      if (pins.length < 3) return null;
      return generatePolygonElement(pins, { fill: polygonFill });
    })
    .filter((p): p is string => p !== null);

  const backgroundRect = includeBackground
    ? `\n  <rect width="${state.imageWidth}" height="${state.imageHeight}" fill="${backgroundColor}"/>`
    : '';

  const polygonsSvg = polygons.length > 0
    ? '\n' + polygons.map(p => `  ${p}`).join('\n')
    : '';

  return `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 ${state.imageWidth} ${state.imageHeight}">${backgroundRect}${polygonsSvg}
</svg>`;
}

// ============================================================================
// Zoom Management
// ============================================================================

export function zoomIn(state: CoordinatePickerState, config: CoordinatePickerConfig = DEFAULT_CONFIG): number {
  if (state.zoomLevel < config.maxZoom) {
    state.zoomLevel = Math.min(state.zoomLevel + config.zoomStep, config.maxZoom);
  }
  return state.zoomLevel;
}

export function zoomOut(state: CoordinatePickerState, config: CoordinatePickerConfig = DEFAULT_CONFIG): number {
  if (state.zoomLevel > config.minZoom) {
    state.zoomLevel = Math.max(state.zoomLevel - config.zoomStep, config.minZoom);
  }
  return state.zoomLevel;
}

export function resetZoom(state: CoordinatePickerState): number {
  state.zoomLevel = 1;
  return state.zoomLevel;
}

// ============================================================================
// State Reset
// ============================================================================

export function clearAllPins(state: CoordinatePickerState): void {
  state.pins = [];
  state.shapes.forEach(shape => {
    shape.pinIds = [];
  });
  state.pinIdCounter = 0;
}

export function resetState(state: CoordinatePickerState, numShapes: number = 4, config: CoordinatePickerConfig = DEFAULT_CONFIG): void {
  const newState = createInitialState(numShapes, config);
  Object.assign(state, newState);
}

// ============================================================================
// Serialization (for saving/loading state)
// ============================================================================

export interface SerializedState {
  pins: Pin[];
  shapes: { color: string; pinIds: number[] }[];
  activeShapeIndex: number;
  pinIdCounter: number;
  imageWidth: number;
  imageHeight: number;
}

export function serializeState(state: CoordinatePickerState): SerializedState {
  return {
    pins: state.pins.map(p => ({ ...p })),
    shapes: state.shapes.map(s => ({ color: s.color, pinIds: [...s.pinIds] })),
    activeShapeIndex: state.activeShapeIndex,
    pinIdCounter: state.pinIdCounter,
    imageWidth: state.imageWidth,
    imageHeight: state.imageHeight,
  };
}

export function deserializeState(serialized: SerializedState): CoordinatePickerState {
  return {
    pins: serialized.pins.map(p => ({ ...p })),
    shapes: serialized.shapes.map(s => ({ color: s.color, pinIds: [...s.pinIds] })),
    activeShapeIndex: serialized.activeShapeIndex,
    pinIdCounter: serialized.pinIdCounter,
    imageWidth: serialized.imageWidth,
    imageHeight: serialized.imageHeight,
    zoomLevel: 1,
  };
}
