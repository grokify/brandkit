import { describe, it, expect, beforeEach } from 'vitest';
import {
  createInitialState,
  createPin,
  deletePin,
  getPinById,
  findNearbyPin,
  addPinToShape,
  removePinFromShape,
  getShapePins,
  getShapesForPin,
  clearShape,
  undoLastPinFromShape,
  setNumShapes,
  isPinUsedByMultipleShapes,
  calculateImageCoordinates,
  calculateDisplayPosition,
  generatePolygonPoints,
  generatePolygonElement,
  generateShapePolygon,
  generateFullSVG,
  zoomIn,
  zoomOut,
  resetZoom,
  clearAllPins,
  resetState,
  serializeState,
  deserializeState,
  CoordinatePickerState,
  DEFAULT_CONFIG,
} from '../src/coordinate-picker';

describe('State Management', () => {
  describe('createInitialState', () => {
    it('creates state with default 4 shapes', () => {
      const state = createInitialState();
      expect(state.shapes.length).toBe(4);
      expect(state.pins).toEqual([]);
      expect(state.activeShapeIndex).toBe(0);
      expect(state.pinIdCounter).toBe(0);
      expect(state.zoomLevel).toBe(1);
    });

    it('creates state with custom number of shapes', () => {
      const state = createInitialState(6);
      expect(state.shapes.length).toBe(6);
    });

    it('assigns colors to shapes cyclically', () => {
      const state = createInitialState(12);
      expect(state.shapes[0].color).toBe(DEFAULT_CONFIG.shapeColors[0]);
      expect(state.shapes[10].color).toBe(DEFAULT_CONFIG.shapeColors[0]); // Cycles back
    });
  });
});

describe('Pin Management', () => {
  let state: CoordinatePickerState;

  beforeEach(() => {
    state = createInitialState(4);
  });

  describe('createPin', () => {
    it('creates a pin with correct properties', () => {
      const pin = createPin(state, 100, 150, 100, 150);
      expect(pin.id).toBe(0);
      expect(pin.x).toBe(100);
      expect(pin.y).toBe(150);
      expect(state.pins.length).toBe(1);
      expect(state.pinIdCounter).toBe(1);
    });

    it('increments pin IDs correctly', () => {
      const pin1 = createPin(state, 10, 20, 10, 20);
      const pin2 = createPin(state, 30, 40, 30, 40);
      expect(pin1.id).toBe(0);
      expect(pin2.id).toBe(1);
      expect(state.pins.length).toBe(2);
    });
  });

  describe('getPinById', () => {
    it('returns pin when found', () => {
      const created = createPin(state, 100, 200, 100, 200);
      const found = getPinById(state, created.id);
      expect(found).toEqual(created);
    });

    it('returns undefined when not found', () => {
      const found = getPinById(state, 999);
      expect(found).toBeUndefined();
    });
  });

  describe('findNearbyPin', () => {
    it('finds pin within threshold', () => {
      const pin = createPin(state, 100, 100, 100, 100);
      const found = findNearbyPin(state.pins, 102, 103);
      expect(found).toEqual(pin);
    });

    it('returns undefined when no pin is close enough', () => {
      createPin(state, 100, 100, 100, 100);
      const found = findNearbyPin(state.pins, 120, 120);
      expect(found).toBeUndefined();
    });
  });

  describe('deletePin', () => {
    it('removes pin from state and all shapes', () => {
      const pin = createPin(state, 100, 100, 100, 100);
      addPinToShape(state, pin.id, 0);
      addPinToShape(state, pin.id, 1);

      deletePin(state, pin.id);

      expect(state.pins.length).toBe(0);
      expect(state.shapes[0].pinIds).not.toContain(pin.id);
      expect(state.shapes[1].pinIds).not.toContain(pin.id);
    });
  });

  describe('isPinUsedByMultipleShapes', () => {
    it('returns false when pin is in one shape', () => {
      const pin = createPin(state, 100, 100, 100, 100);
      addPinToShape(state, pin.id, 0);
      expect(isPinUsedByMultipleShapes(state, pin.id)).toBe(false);
    });

    it('returns true when pin is in multiple shapes', () => {
      const pin = createPin(state, 100, 100, 100, 100);
      addPinToShape(state, pin.id, 0);
      addPinToShape(state, pin.id, 1);
      expect(isPinUsedByMultipleShapes(state, pin.id)).toBe(true);
    });
  });
});

describe('Shape Management', () => {
  let state: CoordinatePickerState;

  beforeEach(() => {
    state = createInitialState(4);
  });

  describe('addPinToShape', () => {
    it('adds pin to shape', () => {
      const pin = createPin(state, 100, 100, 100, 100);
      const result = addPinToShape(state, pin.id, 0);
      expect(result).toBe(true);
      expect(state.shapes[0].pinIds).toContain(pin.id);
    });

    it('returns false for invalid shape index', () => {
      const pin = createPin(state, 100, 100, 100, 100);
      expect(addPinToShape(state, pin.id, -1)).toBe(false);
      expect(addPinToShape(state, pin.id, 10)).toBe(false);
    });

    it('returns false when pin already in shape', () => {
      const pin = createPin(state, 100, 100, 100, 100);
      addPinToShape(state, pin.id, 0);
      expect(addPinToShape(state, pin.id, 0)).toBe(false);
    });
  });

  describe('removePinFromShape', () => {
    it('removes pin from shape', () => {
      const pin = createPin(state, 100, 100, 100, 100);
      addPinToShape(state, pin.id, 0);
      const result = removePinFromShape(state, pin.id, 0);
      expect(result).toBe(true);
      expect(state.shapes[0].pinIds).not.toContain(pin.id);
    });

    it('returns false when pin not in shape', () => {
      const pin = createPin(state, 100, 100, 100, 100);
      expect(removePinFromShape(state, pin.id, 0)).toBe(false);
    });

    it('returns false for invalid shape index', () => {
      const pin = createPin(state, 100, 100, 100, 100);
      expect(removePinFromShape(state, pin.id, -1)).toBe(false);
      expect(removePinFromShape(state, pin.id, 10)).toBe(false);
    });
  });

  describe('getShapePins', () => {
    it('returns pins for shape in order', () => {
      const pin1 = createPin(state, 10, 10, 10, 10);
      const pin2 = createPin(state, 20, 20, 20, 20);
      const pin3 = createPin(state, 30, 30, 30, 30);

      addPinToShape(state, pin1.id, 0);
      addPinToShape(state, pin2.id, 0);
      addPinToShape(state, pin3.id, 0);

      const pins = getShapePins(state, 0);
      expect(pins.length).toBe(3);
      expect(pins[0]).toEqual(pin1);
      expect(pins[1]).toEqual(pin2);
      expect(pins[2]).toEqual(pin3);
    });

    it('returns empty array for invalid shape index', () => {
      expect(getShapePins(state, -1)).toEqual([]);
      expect(getShapePins(state, 10)).toEqual([]);
    });
  });

  describe('clearShape', () => {
    it('clears shape and removes orphaned pins', () => {
      const pin1 = createPin(state, 10, 10, 10, 10);
      const pin2 = createPin(state, 20, 20, 20, 20);

      addPinToShape(state, pin1.id, 0);
      addPinToShape(state, pin2.id, 0);
      addPinToShape(state, pin2.id, 1); // pin2 used by another shape

      const removed = clearShape(state, 0);

      expect(state.shapes[0].pinIds).toEqual([]);
      expect(removed).toContain(pin1.id); // Orphaned, should be removed
      expect(removed).not.toContain(pin2.id); // Still used
      expect(state.pins.find(p => p.id === pin1.id)).toBeUndefined();
      expect(state.pins.find(p => p.id === pin2.id)).toBeDefined();
    });

    it('returns empty array for invalid shape index', () => {
      expect(clearShape(state, -1)).toEqual([]);
      expect(clearShape(state, 10)).toEqual([]);
    });
  });

  describe('undoLastPinFromShape', () => {
    it('removes last pin and deletes if orphaned', () => {
      const pin = createPin(state, 100, 100, 100, 100);
      addPinToShape(state, pin.id, 0);

      const result = undoLastPinFromShape(state, 0);

      expect(result).toEqual({ pinId: pin.id, deleted: true });
      expect(state.shapes[0].pinIds).toEqual([]);
      expect(state.pins.length).toBe(0);
    });

    it('removes last pin but keeps if still used', () => {
      const pin = createPin(state, 100, 100, 100, 100);
      addPinToShape(state, pin.id, 0);
      addPinToShape(state, pin.id, 1);

      const result = undoLastPinFromShape(state, 0);

      expect(result).toEqual({ pinId: pin.id, deleted: false });
      expect(state.pins.length).toBe(1);
    });

    it('returns null when shape is empty', () => {
      expect(undoLastPinFromShape(state, 0)).toBeNull();
    });

    it('returns null for invalid shape index', () => {
      expect(undoLastPinFromShape(state, -1)).toBeNull();
      expect(undoLastPinFromShape(state, 10)).toBeNull();
    });
  });

  describe('setNumShapes', () => {
    it('adds new shapes', () => {
      setNumShapes(state, 6);
      expect(state.shapes.length).toBe(6);
    });

    it('preserves existing shapes when reducing', () => {
      const pin = createPin(state, 100, 100, 100, 100);
      addPinToShape(state, pin.id, 0);

      setNumShapes(state, 2);

      expect(state.shapes.length).toBe(2);
      expect(state.shapes[0].pinIds).toContain(pin.id);
    });

    it('adjusts activeShapeIndex when needed', () => {
      state.activeShapeIndex = 3;
      setNumShapes(state, 2);
      expect(state.activeShapeIndex).toBe(0);
    });
  });
});

describe('Coordinate Calculations', () => {
  describe('calculateImageCoordinates', () => {
    it('converts click coords at zoom 1', () => {
      const result = calculateImageCoordinates(100, 150, 1);
      expect(result).toEqual({ x: 100, y: 150 });
    });

    it('converts click coords at zoom 2', () => {
      const result = calculateImageCoordinates(200, 300, 2);
      expect(result).toEqual({ x: 100, y: 150 });
    });

    it('rounds coordinates', () => {
      const result = calculateImageCoordinates(100, 150, 1.5);
      expect(result.x).toBe(67);
      expect(result.y).toBe(100);
    });
  });

  describe('calculateDisplayPosition', () => {
    it('converts image coords to display at zoom 1', () => {
      const result = calculateDisplayPosition(100, 150, 1);
      expect(result).toEqual({ x: 100, y: 150 });
    });

    it('converts image coords to display at zoom 2', () => {
      const result = calculateDisplayPosition(100, 150, 2);
      expect(result).toEqual({ x: 200, y: 300 });
    });
  });
});

describe('SVG Generation', () => {
  let state: CoordinatePickerState;

  beforeEach(() => {
    state = createInitialState(2);
    state.imageWidth = 400;
    state.imageHeight = 400;
  });

  describe('generatePolygonPoints', () => {
    it('generates point string', () => {
      const pins = [
        { id: 0, x: 100, y: 100, displayX: 100, displayY: 100 },
        { id: 1, x: 200, y: 100, displayX: 200, displayY: 100 },
        { id: 2, x: 150, y: 200, displayX: 150, displayY: 200 },
      ];
      expect(generatePolygonPoints(pins)).toBe('100,100 200,100 150,200');
    });
  });

  describe('generatePolygonElement', () => {
    const pins = [
      { id: 0, x: 100, y: 100, displayX: 100, displayY: 100 },
      { id: 1, x: 200, y: 100, displayX: 200, displayY: 100 },
      { id: 2, x: 150, y: 200, displayX: 150, displayY: 200 },
    ];

    it('returns empty string for less than 3 pins', () => {
      expect(generatePolygonElement(pins.slice(0, 2))).toBe('');
    });

    it('generates polygon with default fill', () => {
      const svg = generatePolygonElement(pins);
      expect(svg).toContain('points="100,100 200,100 150,200"');
      expect(svg).toContain('fill="#00FF00"');
    });

    it('generates polygon with custom options', () => {
      const svg = generatePolygonElement(pins, {
        fill: '#FF0000',
        stroke: '#000',
        strokeWidth: 2,
      });
      expect(svg).toContain('fill="#FF0000"');
      expect(svg).toContain('stroke="#000"');
      expect(svg).toContain('stroke-width="2"');
    });
  });

  describe('generateFullSVG', () => {
    it('generates SVG with background', () => {
      const pin1 = createPin(state, 100, 100, 100, 100);
      const pin2 = createPin(state, 200, 100, 200, 100);
      const pin3 = createPin(state, 150, 200, 150, 200);

      addPinToShape(state, pin1.id, 0);
      addPinToShape(state, pin2.id, 0);
      addPinToShape(state, pin3.id, 0);

      const svg = generateFullSVG(state);

      expect(svg).toContain('xmlns="http://www.w3.org/2000/svg"');
      expect(svg).toContain('viewBox="0 0 400 400"');
      expect(svg).toContain('<rect');
      expect(svg).toContain('fill="#000000"');
      expect(svg).toContain('<polygon');
    });

    it('generates SVG without background', () => {
      const pin1 = createPin(state, 100, 100, 100, 100);
      const pin2 = createPin(state, 200, 100, 200, 100);
      const pin3 = createPin(state, 150, 200, 150, 200);

      addPinToShape(state, pin1.id, 0);
      addPinToShape(state, pin2.id, 0);
      addPinToShape(state, pin3.id, 0);

      const svg = generateFullSVG(state, { includeBackground: false });

      expect(svg).not.toContain('<rect');
    });

    it('skips shapes with less than 3 pins', () => {
      const pin1 = createPin(state, 100, 100, 100, 100);
      const pin2 = createPin(state, 200, 100, 200, 100);

      addPinToShape(state, pin1.id, 0);
      addPinToShape(state, pin2.id, 0);

      const svg = generateFullSVG(state);

      expect(svg).not.toContain('<polygon');
    });
  });
});

describe('Zoom Management', () => {
  let state: CoordinatePickerState;

  beforeEach(() => {
    state = createInitialState();
  });

  describe('zoomIn', () => {
    it('increases zoom by step', () => {
      const result = zoomIn(state);
      expect(result).toBe(1.25);
      expect(state.zoomLevel).toBe(1.25);
    });

    it('respects max zoom', () => {
      state.zoomLevel = 3.9;
      const result = zoomIn(state);
      expect(result).toBe(4);
    });

    it('does not exceed max zoom', () => {
      state.zoomLevel = 4;
      const result = zoomIn(state);
      expect(result).toBe(4);
    });
  });

  describe('zoomOut', () => {
    it('decreases zoom by step', () => {
      const result = zoomOut(state);
      expect(result).toBe(0.75);
      expect(state.zoomLevel).toBe(0.75);
    });

    it('respects min zoom', () => {
      state.zoomLevel = 0.6;
      const result = zoomOut(state);
      expect(result).toBe(0.5);
    });

    it('does not go below min zoom', () => {
      state.zoomLevel = 0.5;
      const result = zoomOut(state);
      expect(result).toBe(0.5);
    });
  });

  describe('resetZoom', () => {
    it('resets zoom to 1', () => {
      state.zoomLevel = 2.5;
      const result = resetZoom(state);
      expect(result).toBe(1);
      expect(state.zoomLevel).toBe(1);
    });
  });
});

describe('State Reset', () => {
  let state: CoordinatePickerState;

  beforeEach(() => {
    state = createInitialState();
    // Add some data
    const pin = createPin(state, 100, 100, 100, 100);
    addPinToShape(state, pin.id, 0);
  });

  describe('clearAllPins', () => {
    it('clears all pins and shape references', () => {
      clearAllPins(state);
      expect(state.pins).toEqual([]);
      expect(state.shapes.every(s => s.pinIds.length === 0)).toBe(true);
      expect(state.pinIdCounter).toBe(0);
    });
  });
});

describe('Serialization', () => {
  let state: CoordinatePickerState;

  beforeEach(() => {
    state = createInitialState(2);
    state.imageWidth = 300;
    state.imageHeight = 300;

    const pin1 = createPin(state, 10, 20, 10, 20);
    const pin2 = createPin(state, 30, 40, 30, 40);
    addPinToShape(state, pin1.id, 0);
    addPinToShape(state, pin2.id, 0);
    addPinToShape(state, pin2.id, 1); // Shared pin
  });

  describe('serializeState', () => {
    it('creates serializable object', () => {
      const serialized = serializeState(state);

      expect(serialized.pins.length).toBe(2);
      expect(serialized.shapes.length).toBe(2);
      expect(serialized.activeShapeIndex).toBe(0);
      expect(serialized.pinIdCounter).toBe(2);
      expect(serialized.imageWidth).toBe(300);
      expect(serialized.imageHeight).toBe(300);
    });

    it('creates deep copies', () => {
      const serialized = serializeState(state);

      // Modify original
      state.pins[0].x = 999;
      state.shapes[0].pinIds.push(999);

      // Serialized should be unchanged
      expect(serialized.pins[0].x).toBe(10);
      expect(serialized.shapes[0].pinIds.length).toBe(2);
    });
  });

  describe('deserializeState', () => {
    it('restores state from serialized data', () => {
      const serialized = serializeState(state);
      const restored = deserializeState(serialized);

      expect(restored.pins.length).toBe(2);
      expect(restored.shapes.length).toBe(2);
      expect(restored.activeShapeIndex).toBe(0);
      expect(restored.pinIdCounter).toBe(2);
      expect(restored.imageWidth).toBe(300);
      expect(restored.imageHeight).toBe(300);
      expect(restored.zoomLevel).toBe(1); // Reset to 1
    });

    it('creates deep copies', () => {
      const serialized = serializeState(state);
      const restored = deserializeState(serialized);

      // Modify serialized
      serialized.pins[0].x = 999;

      // Restored should be unchanged
      expect(restored.pins[0].x).toBe(10);
    });
  });

  it('roundtrip preserves data', () => {
    const serialized = serializeState(state);
    const restored = deserializeState(serialized);
    const reSerialized = serializeState(restored);

    expect(reSerialized).toEqual(serialized);
  });
});

describe('getShapesForPin', () => {
  let state: CoordinatePickerState;

  beforeEach(() => {
    state = createInitialState(4);
  });

  it('returns empty array when pin is not in any shape', () => {
    const pin = createPin(state, 100, 100, 100, 100);
    const shapes = getShapesForPin(state, pin.id);
    expect(shapes).toEqual([]);
  });

  it('returns single shape when pin is in one shape', () => {
    const pin = createPin(state, 100, 100, 100, 100);
    addPinToShape(state, pin.id, 1);

    const shapes = getShapesForPin(state, pin.id);

    expect(shapes.length).toBe(1);
    expect(shapes[0]).toBe(state.shapes[1]);
  });

  it('returns multiple shapes when pin is shared', () => {
    const pin = createPin(state, 100, 100, 100, 100);
    addPinToShape(state, pin.id, 0);
    addPinToShape(state, pin.id, 2);
    addPinToShape(state, pin.id, 3);

    const shapes = getShapesForPin(state, pin.id);

    expect(shapes.length).toBe(3);
    expect(shapes).toContain(state.shapes[0]);
    expect(shapes).toContain(state.shapes[2]);
    expect(shapes).toContain(state.shapes[3]);
  });
});

describe('generateShapePolygon', () => {
  let state: CoordinatePickerState;

  beforeEach(() => {
    state = createInitialState(2);
    state.imageWidth = 400;
    state.imageHeight = 400;
  });

  it('returns empty string when shape has less than 3 pins', () => {
    const pin1 = createPin(state, 100, 100, 100, 100);
    const pin2 = createPin(state, 200, 100, 200, 100);

    addPinToShape(state, pin1.id, 0);
    addPinToShape(state, pin2.id, 0);

    const svg = generateShapePolygon(state, 0);
    expect(svg).toBe('');
  });

  it('generates polygon for shape with 3+ pins', () => {
    const pin1 = createPin(state, 100, 100, 100, 100);
    const pin2 = createPin(state, 200, 100, 200, 100);
    const pin3 = createPin(state, 150, 200, 150, 200);

    addPinToShape(state, pin1.id, 0);
    addPinToShape(state, pin2.id, 0);
    addPinToShape(state, pin3.id, 0);

    const svg = generateShapePolygon(state, 0);

    expect(svg).toContain('<polygon');
    expect(svg).toContain('points="100,100 200,100 150,200"');
    expect(svg).toContain('fill="#00FF00"');
  });

  it('generates polygon with custom options', () => {
    const pin1 = createPin(state, 10, 10, 10, 10);
    const pin2 = createPin(state, 20, 10, 20, 10);
    const pin3 = createPin(state, 15, 20, 15, 20);

    addPinToShape(state, pin1.id, 1);
    addPinToShape(state, pin2.id, 1);
    addPinToShape(state, pin3.id, 1);

    const svg = generateShapePolygon(state, 1, { fill: '#FF0000', stroke: '#000' });

    expect(svg).toContain('fill="#FF0000"');
    expect(svg).toContain('stroke="#000"');
  });

  it('returns empty string for invalid shape index', () => {
    expect(generateShapePolygon(state, -1)).toBe('');
    expect(generateShapePolygon(state, 10)).toBe('');
  });
});

describe('resetState', () => {
  it('resets state to initial values', () => {
    const state = createInitialState(4);

    // Add some data
    const pin1 = createPin(state, 100, 100, 100, 100);
    const pin2 = createPin(state, 200, 200, 200, 200);
    addPinToShape(state, pin1.id, 0);
    addPinToShape(state, pin2.id, 1);
    state.activeShapeIndex = 2;
    state.zoomLevel = 2.5;
    state.imageWidth = 800;
    state.imageHeight = 600;

    // Reset
    resetState(state);

    // Verify reset
    expect(state.pins).toEqual([]);
    expect(state.pinIdCounter).toBe(0);
    expect(state.activeShapeIndex).toBe(0);
    expect(state.zoomLevel).toBe(1);
    expect(state.imageWidth).toBe(400);
    expect(state.imageHeight).toBe(400);
    expect(state.shapes.length).toBe(4);
    expect(state.shapes.every(s => s.pinIds.length === 0)).toBe(true);
  });

  it('resets state with custom number of shapes', () => {
    const state = createInitialState(2);

    // Add data
    createPin(state, 50, 50, 50, 50);

    // Reset with different shape count
    resetState(state, 6);

    expect(state.shapes.length).toBe(6);
    expect(state.pins).toEqual([]);
  });
});
