<template>
  <div class="test-settings-page">
    <h2>测试插件·自定义设置页</h2>
    <p>本页经 plugin.json 根级 settingsPage 声明挂载（vueSource 运行时编译形态），插件管理页「设置」直达此处；模板内 el-button 经宿主全局注册解析，反馈走 __PLUGIN_CTX__.globals 契约通道。</p>
    <el-button type="primary" @click="handleClick">点我反馈</el-button>
    <p v-if="lastClickedAt !== ''">按钮已点击：{{ lastClickedAt }}</p>
  </div>
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
.test-settings-page {
  padding: 20px;
  color: var(--app-text-primary);
}
</style>
