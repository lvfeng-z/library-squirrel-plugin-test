// 预编译工厂函数格式：export default function(__VUE__, __WAILS_RUNTIME__) { return defineComponent({...}) }
// 主程序 loadCompiledComponent 调用 module.default(Vue, WailsRuntime) 注入依赖，避免插件自带 vue/wails 依赖。
// 禁止用 import { X as Y } 的 as 语法（Vite 工厂插件不兼容）——需要别名直接在解构时改名。
export default function (__VUE__, __WAILS_RUNTIME__) {
  const { computed, h, defineComponent } = __VUE__
  return defineComponent({
    name: 'TestArticleViewer',
    props: {
      resource: { type: Object, default: () => ({}) },
      work: { type: Object, default: () => ({}) }
    },
    setup(props) {
      // 展示注入的 resource/work 关键字段，证明插件渲染器契约（{resource, work}）生效
      const info = computed(
        () =>
          JSON.stringify(
            {
              resourceType: props.resource?.resourceType,
              resourceId: props.resource?.id,
              storeCount: props.resource?.stores?.length ?? 0,
              workName: props.work?.work?.siteWorkName
            },
            null,
            2
          )
      )
      return () =>
        h('div', { class: 'test-article-viewer' }, [
          h('h2', { class: 'test-article-viewer-title' }, '🔧 测试插件渲染器（article）'),
          h(
            'p',
            { class: 'test-article-viewer-desc' },
            '此区域由测试插件的 resourceViewer 渲染，已覆盖内置 ArticleRenderer——若看到本提示，说明插件渲染器 Handler 通道端到端打通。'
          ),
          h('pre', { class: 'test-article-viewer-info' }, info.value)
        ])
    }
  })
}
