<script setup>
import ClaudeConfigManager from './components/ClaudeConfigManager.vue'
</script>

<template>
  <div class="app">
    <div class="background-decoration">
      <div class="decoration-circle circle-1"></div>
      <div class="decoration-circle circle-2"></div>
      <div class="decoration-circle circle-3"></div>
    </div>
    
    <header class="app-header">
      <div class="logo-container">
        <div class="logo-glow"></div>
      </div>
      <h1>Claude Code 配置管理器</h1>
      <p class="subtitle">优雅管理您的 AI 助手配置</p>
    </header>
    
    <main class="app-main">
      <ClaudeConfigManager />
    </main>
    
    <footer class="app-footer">
      <p>Made with ❤️ for Claude Code</p>
    </footer>
  </div>
</template>

<style>
* {
  box-sizing: border-box;
  margin: 0;
  padding: 0;
}

body {
  font-family: -apple-system, BlinkMacSystemFont, 'SF Pro Display', 'SF Pro Text', 'Helvetica Neue', Helvetica, Arial, sans-serif;
  background: #f9f9f9;
  min-height: 100vh;
  overflow-x: hidden;
}

.app {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 2rem 1rem;
  position: relative;
  background: linear-gradient(135deg, #f9f9f9 0%, #f0f0f0 50%, #e8e8e8 100%);
}

.background-decoration {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  pointer-events: none;
  z-index: 0;
}

.decoration-circle {
  position: absolute;
  border-radius: 50%;
  background: linear-gradient(45deg, rgba(0, 122, 255, 0.08), rgba(52, 199, 89, 0.06));
  backdrop-filter: blur(2px);
  animation: float 6s ease-in-out infinite;
}

.circle-1 {
  width: 300px;
  height: 300px;
  top: 10%;
  left: -10%;
  animation-delay: 0s;
}

.circle-2 {
  width: 200px;
  height: 200px;
  top: 60%;
  right: -5%;
  animation-delay: 2s;
}

.circle-3 {
  width: 150px;
  height: 150px;
  bottom: 20%;
  left: 20%;
  animation-delay: 4s;
}

@keyframes float {
  0%, 100% { transform: translateY(0px) rotate(0deg); }
  50% { transform: translateY(-20px) rotate(180deg); }
}

.app-header {
  text-align: center;
  margin-bottom: 3rem;
  color: #1c1c1e;
  z-index: 1;
  position: relative;
}

.logo-container {
  position: relative;
  display: inline-block;
  margin-bottom: 1.5rem;
}

#logo {
  width: 100px;
  height: 100px;
  border-radius: 24px;
  box-shadow: 0 20px 40px rgba(0, 0, 0, 0.3);
  transition: transform 0.3s ease;
  position: relative;
  z-index: 2;
}

.logo-glow {
  position: absolute;
  top: -10px;
  left: -10px;
  right: -10px;
  bottom: -10px;
  background: linear-gradient(45deg, #007AFF, #34C759, #5856D6);
  border-radius: 34px;
  opacity: 0.3;
  filter: blur(20px);
  z-index: 1;
  animation: pulse 2s ease-in-out infinite;
}

@keyframes pulse {
  0%, 100% { opacity: 0.3; transform: scale(1); }
  50% { opacity: 0.5; transform: scale(1.05); }
}

#logo:hover {
  transform: scale(1.05) rotate(5deg);
}

.app-header h1 {
  font-size: 3rem;
  font-weight: 700;
  margin: 0.5rem 0;
  background: linear-gradient(45deg, #1c1c1e, #48484a);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
  text-shadow: none;
  letter-spacing: -0.02em;
}

.subtitle {
  font-size: 1.2rem;
  opacity: 0.7;
  margin: 0;
  font-weight: 400;
  letter-spacing: 0.5px;
  color: #8e8e93;
}

.app-main {
  width: 100%;
  max-width: 500px;
  flex: 1;
  display: flex;
  align-items: flex-start;
  z-index: 1;
  position: relative;
}

.app-footer {
  margin-top: 2rem;
  color: #8e8e93;
  font-size: 0.9rem;
  z-index: 1;
  position: relative;
}

/* Responsive design */
@media (max-width: 768px) {
  .app {
    padding: 1.5rem 1rem;
  }
  
  .app-header h1 {
    font-size: 2.2rem;
  }
  
  .subtitle {
    font-size: 1rem;
  }
  
  #logo {
    width: 80px;
    height: 80px;
  }
}

@media (max-width: 480px) {
  .app-header h1 {
    font-size: 1.8rem;
  }
  
  .circle-1, .circle-2, .circle-3 {
    display: none;
  }
}
</style>
