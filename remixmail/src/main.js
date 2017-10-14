// The Vue build version to load with the `import` command
// (runtime-only or standalone) has been set in webpack.base.conf with an alias.
import Vue from 'vue'
import App from './App'
import router from './router'
import VueResource from 'vue-resource'
import VueMaterial from 'vue-material'
import 'vue-material/dist/vue-material.css'

Vue.use(VueResource)
Vue.use(VueMaterial)

Vue.config.productionTip = false

Vue.material.registerTheme({
  default: {
    primary: {
      color: 'black',
      hue: 700,
      background: 'blue'
    },
    accent: 'red'
  }
})

/* eslint-disable no-new */
new Vue({
  el: '#app',
  router,
  template: '<App/>',
  components: { App }
})
