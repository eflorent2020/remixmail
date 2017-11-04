import Vue from 'vue'
import Router from 'vue-router'
import HelloWorld from '@/components/HelloWorld'
import Validator from '@/components/Validator'
import Admin from '@/components/Admin'
// import App from '@/App'

Vue.use(Router)

export default new Router({
  routes: [
    {
      path: '/',
      name: 'Hello',
      component: HelloWorld,
      props: true
    },
    {
      path: '/alias/validate/:validationKey',
      name: 'Validator',
      component: Validator
    },
    {
      path: '/admin',
      name: 'Admin',
      component: Admin,
      props: true
    }
  ]
})
