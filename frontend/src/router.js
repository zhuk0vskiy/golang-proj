import {createWebHistory, createRouter} from 'vue-router';

import Login from './components/Login.vue';
import Login from './components/Register.vue';

// const Admin = () => import('./components/Admin.vue');

const routes = [
//   {
//     path: '/legacy',
//     name: 'home',
//     component: Home,
//   },
//   {
//     path: '/legacy/home',
//     component: Home,
//   },
  {
    path: '/legacy/login',
    component: Login,
  },
  {
    path: '/legacy/register',
    component: Register,
  },
//   {
//     path: '/legacy/profile',
//     component: Employee,
//   },
//   {
//     path: '/legacy/find-employees',
//     component: Admin,
//   },
];

const router = createRouter({
  history: createWebHistory(),
  routes: routes,
});

router.beforeEach((to, from, next) => {
  const publicPages = ['/legacy/login', '/legacy/register'];
  const authRequired = !publicPages.includes(to.path);
  const loggedIn = localStorage.getItem('user');

  if (authRequired && !loggedIn) {
    next('/legacy/login');
  } else {
    next();
  }
});

export default router;