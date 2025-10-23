# Firefly Game Boilerplate Development Plan

## 🎯 HIGH PRIORITY - Core Game Mechanics

### Technical debits and wishlist

- Create shape package to split some things from physics package.
- Review player platform and top down. Maybe it should be removed.
- Add cam movement command to sequence
- Change physics in a real time sandbox control
- reduce repeated scene contents
- change sprite animation when jumping

### 4. UI/HUD System

- Button system with click handling

## 🔧 MEDIUM PRIORITY - Game Systems

### 5. Scoring & Progression System

- Points system with high scores
- ~Level progression mechanics~
- Achievement system

### 6. Collectibles System

- Interactive items (coins, power-ups)
- Item effects and behaviors
- Collection feedback and animations

### 7. Enhanced Animation System

- Sprite interpolation and easing
- Complex animation sequences
- Animation state management

### 8. Camera System

- ~Camera following and zoom~
- World scrolling for larger maps
- Viewport management

## 🎨 LOWER PRIORITY - Polish & Advanced Features

### 9. Enhanced Input System

- Gamepad/controller support
- Key remapping and configuration
- Input event system

### 10. Improved Collision System

- Different collision types (solid, trigger, sensor)
- Collision layers and filtering
- Event-driven collision responses

### 11. Particle Effects System

- Particle emitters and effects
- Visual feedback for actions
- Performance optimization

### 12. Level/Map System

- Tilemap support and loading
- Level editor integration
- Dynamic level generation

### 13. Save System

- Game state persistence
- Settings and configuration saving
- Save file management

### 14. Performance Optimization

- Rendering optimization
- Memory management
- Profiling and debugging tools

## 🏗️ FUTURE - Framework Conversion (Not Priority)

### 15. Framework Core Abstraction

- Create reusable framework interfaces
- Entity-component system
- Plugin architecture

### 16. Framework Systems

- Asset management system
- Event/messaging system
- Configuration management

### 17. Developer Experience

- Public API design
- Documentation and tutorials
- Example games and templates

### 18. Framework Distribution

- Version management strategy
- Update workflow design
- Community tools and support

---

## Current Status

✅ Basic player movement and physics  
✅ Collision detection system  
✅ Sprite animation (idle/walk)  
✅ Input handling (keyboard)  
✅ Boundary and obstacle system

## Next Focus

🎯 **Start with Game State Management**

- This provides the foundation for all other systems and creates a proper game structure.
