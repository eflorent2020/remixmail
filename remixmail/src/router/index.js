import Vue from 'vue'
import Router from 'vue-router'
import HelloWorld from '@/components/HelloWorld'
import Validator from '@/components/Validator'
import Admin from '@/components/Admin'

Vue.use(Router)

export default new Router({
  routes: [
    {
      path: '/',
      name: 'Hello',
      component: HelloWorld
    },
    {
      path: '/alias/validate/:validationKey',
      name: 'Validator',
      component: Validator
    },
    {
      path: '/admin',
      name: 'Admin',
      component: Admin
    }
  ]
})
