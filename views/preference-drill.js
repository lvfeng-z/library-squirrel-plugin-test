// 预编译工厂函数格式：export default function(__VUE__, __WAILS_RUNTIME__) { return defineComponent({...}) }
// 主程序 loadCompiledComponent 调用 module.default(Vue, WailsRuntime) 注入依赖，避免插件自带 vue/wails 依赖。
// 禁止用 import { X as Y } 的 as 语法（Vite 工厂插件不兼容）——需要别名直接在解构时改名。
//
// 偏好三方法演练页：经 Wails 事件（plugin-dev-guide 第九章前后端通信）向插件 Go 侧发
// 演练请求（写入/覆写/读取/列键），Go 侧经 SDK 偏好域三方法执行并回发结果，本页展示
// 执行结果——「插件经 SDK 写入/读取/列键」的真实交互面。
export default function (__VUE__, __WAILS_RUNTIME__) {
  const { ref, h, defineComponent, onMounted, onUnmounted } = __VUE__
  const REQ_TOPIC = 'plugin:test-plugin:preference:request'
  const RESP_TOPIC = 'plugin:test-plugin:preference:response'
  const DRILL_KEY = 'download.form'

  // 插件 Go 侧 []byte 载荷经 Wails 事件序列化到达（JSON 文本或 base64 文本），统一解出对象
  function decodePayload (raw) {
    if (raw == null) return null
    if (typeof raw === 'object') return raw
    if (typeof raw === 'string') {
      try { return JSON.parse(raw) } catch (e) { /* 尝试 base64 解码 */ }
      try { return JSON.parse(atob(raw)) } catch (e) { /* 放弃 */ }
    }
    return null
  }

  return defineComponent({
    name: 'TestPreferenceDrill',
    setup () {
      const seq = ref(0)
      const pending = ref(null) // { seq, label } 待配对请求
      const logLines = ref([]) // 演练日志行（新在前）
      let offResp = null
      let timer = null
      const events = __WAILS_RUNTIME__ && __WAILS_RUNTIME__.Events

      function pushLine (text, kind) {
        logLines.value.unshift({ text: text, kind: kind || 'info' })
        if (logLines.value.length > 30) logLines.value.pop()
      }

      // 结果行渲染：get 区分有记录/无记录，list 展示键清单，写类展示 title
      function renderResult (label, p) {
        if (!p.ok) return label + ' ✗ ' + (p.error || '未知错误')
        if (p.action === 'get') {
          return p.found
            ? label + ' ✓ 有记录 title=' + p.title + '｜data=' + p.data
            : label + ' ✓ 无记录（下次问答将重新发起）'
        }
        if (p.action === 'list') return label + ' ✓ 键清单 [' + (p.keys || []).join(', ') + ']'
        return label + ' ✓ 已写入 title=' + p.title
      }

      function onResp (event) {
        const payload = decodePayload(event && event.data)
        if (!payload || pending.value == null || payload.seq !== pending.value.seq) return
        clearTimeout(timer)
        pushLine(renderResult(pending.value.label, payload), payload.ok ? 'ok' : 'err')
        pending.value = null
      }

      onMounted(function () {
        if (!events) {
          pushLine('Wails Runtime 事件通道不可用，演练无法进行', 'err')
          return
        }
        offResp = events.On(RESP_TOPIC, onResp)
        pushLine('演练页已就绪（请求 topic=' + REQ_TOPIC + '）')
      })
      onUnmounted(function () {
        if (offResp) offResp()
        clearTimeout(timer)
      })

      function drill (action, label) {
        if (!events) return
        seq.value += 1
        pending.value = { seq: seq.value, label: label }
        pushLine('→ ' + label + '（action=' + action + ' seq=' + seq.value + '）')
        events.Emit(REQ_TOPIC, JSON.stringify({ action: action, seq: seq.value }))
        clearTimeout(timer)
        timer = setTimeout(function () {
          if (pending.value && pending.value.seq === seq.value) {
            pushLine(label + ' ✗ 8s 内未收到插件回发（插件进程未运行或通道异常）', 'err')
            pending.value = null
          }
        }, 8000)
      }

      return function () {
        return h('div', { class: 'pref-drill' }, [
          h('h2', { class: 'pref-drill-title' }, '🔧 测试插件·偏好三方法演练'),
          h('p', { class: 'pref-drill-desc' }, [
            '经 SDK 偏好域（GetPreference/SetPreference/ListMyPreferences）演练写入、覆写、读取与列键；固定键 ',
            h('code', null, DRILL_KEY),
            '；结果同步展示在主程序「设置 → 记住的选择 → 插件偏好」与本插件设置对话框。'
          ]),
          h('div', { class: 'pref-drill-actions' }, [
            h('button', {
              class: 'pref-drill-btn pref-drill-btn-set',
              onClick: function () { drill('set', '写入偏好（图文形态）') }
            }, '写入偏好（图文形态）'),
            h('button', {
              class: 'pref-drill-btn pref-drill-btn-overwrite',
              onClick: function () { drill('overwrite', '覆写偏好（文章模式）') }
            }, '覆写偏好（文章模式）'),
            h('button', {
              class: 'pref-drill-btn pref-drill-btn-get',
              onClick: function () { drill('get', '读取偏好') }
            }, '读取偏好'),
            h('button', {
              class: 'pref-drill-btn pref-drill-btn-list',
              onClick: function () { drill('list', '列键清单') }
            }, '列键清单')
          ]),
          h('div', { class: 'pref-drill-log' }, logLines.value.map(function (line, i) {
            return h('div', {
              class: 'pref-drill-log-line pref-drill-log-' + line.kind + (i === 0 ? ' pref-drill-log-latest' : ''),
              'data-drill-latest': i === 0 ? '1' : null
            }, line.text)
          }))
        ])
      }
    }
  })
}
