// 预编译工厂函数格式：export default function(__VUE__, __WAILS_RUNTIME__) { return defineComponent({...}) }
// 主程序 loadCompiledComponent 调用 module.default(Vue, WailsRuntime) 注入依赖，避免插件自带 vue/wails 依赖。
// 禁止用 import { X as Y } 的 as 语法（Vite 工厂插件不兼容）——需要别名直接在解构时改名。
//
// 偏好演练子页：plugin.json 中经 content.parent 声明为偏好演练页（test-preference-drill）
// 的子页面，本身无业务交互——供宿主面包屑、Alt+← 上退与停用级联落点的链路验证，
// 页面标题文案即断言锚点。
export default function (__VUE__, __WAILS_RUNTIME__) {
  const { ref, h, defineComponent, onMounted, onUnmounted } = __VUE__

  return defineComponent({
    name: 'TestPreferenceDrillDetail',
    setup () {
      const routeInfo = ref('')

      let offHash = null
      function syncRouteInfo () {
        routeInfo.value = location.hash || '(空)'
      }

      onMounted(function () {
        syncRouteInfo()
        window.addEventListener('hashchange', syncRouteInfo)
        offHash = function () { window.removeEventListener('hashchange', syncRouteInfo) }
      })
      onUnmounted(function () {
        if (offHash) offHash()
      })

      return function () {
        return h('div', { class: 'pref-drill-detail' }, [
          h('h2', { class: 'pref-drill-detail-title' }, '🧩 测试插件·偏好演练子页'),
          h('p', { class: 'pref-drill-detail-desc' }, [
            '本页经清单 content.parent 声明为「测试·偏好演练」的子页面：顶部面包屑显示祖先链，',
            h('br'),
            'Alt+← 逐级上退回父页面；插件停用/卸载时与父页同批消亡，落点为第一个幸存祖先。'
          ]),
          h('p', { class: 'pref-drill-detail-route' }, [
            '当前地址：',
            h('code', null, routeInfo.value)
          ])
        ])
      }
    }
  })
}
