<template>
  <BaseView>
    <div class="test-settings-container">
      <el-scrollbar>
        <div class="test-settings">
          <h2>测试插件·自定义设置页</h2>
          <p>本页经 plugin.json 根级 settingsPage 声明挂载（vueSource 运行时编译形态），插件管理页「设置」直达此处；页骨架即宿主标准形态——BaseView 外壳 + 内容容器 10px 边距 + 内部 el-scrollbar 承滚动，与 el-button 同为零 import 的宿主全局注册组件，反馈走 __PLUGIN_CTX__.globals 契约通道。</p>
          <el-button type="primary" @click="handleClick">点我反馈</el-button>
          <p v-if="lastClickedAt !== ''">按钮已点击：{{ lastClickedAt }}</p>
        </div>
      </el-scrollbar>
    </div>
  </BaseView>
</template>

<script>
export default {
  name: 'TestSettingsPage',
  data() {
    return { lastClickedAt: '' }
  },
  methods: {
    handleClick() {
      this.lastClickedAt = new Date().toLocaleTimeString()
      window.__PLUGIN_CTX__.globals.$message.success('自定义设置页按钮已点击')
    }
  }
}
</script>

<style>
.test-settings-container {
  height: calc(100% - 20px);
  width: calc(100% - 20px);
  margin: 10px;
}

.test-settings {
  width: 100%;
  color: var(--app-text-primary);
}
</style>
