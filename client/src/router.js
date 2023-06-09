import { createRouter, createWebHashHistory } from 'vue-router';
import TheHome from './components/TheHome.vue';
import TheRegistration from './components/TheRegistration.vue';
import UserDash from './components/UserDash.vue';

const router=createRouter({
    history:createWebHashHistory(),
    routes:[
        {path:"/",redirect:'/home'},

        {path:"/home",component:TheHome},

        {path:"/registration",component:TheRegistration},

        {path:"/user",component:UserDash},
    ]
})

export default router;