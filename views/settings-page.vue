<template>
  <!--
    弹窗内容形态：本组件是被动内容，由宿主设置弹窗壳动态挂载。
    根元素即内容容器，按自然流高度撑开，不使用整页骨架，也不自带滚动容器
    （高度封顶与滚动条由宿主壳的滚动容器负责）。
  -->
  <div class="test-settings">
    <h2>测试插件·自定义设置页</h2>
    <p>本页经 plugin.json 根级 settingsPage 声明挂载（vueSource 运行时编译形态）；清单根级 settingsPresent 声明为 dialog 后，插件管理页「设置」在宿主弹窗内容区呈现本组件，页面骨架与滚动纪律归宿主壳。组件内只用零 import 的宿主全局注册组件（el-button），反馈经 __PLUGIN_CTX__.globals 契约通道。</p>
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
.test-settings {
  width: 100%;
  color: var(--app-text-primary);
}
</style>
