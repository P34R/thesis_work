export default{
  
    namespaced:true,
    state(){
        return{
            conn:null,
            loggedUsername: "",
            loggedUserPrivateKey:null,
            //chatPrivateKeys:new Map(),
            //ChatterName:"",
        };
    },
    mutations:{
        setUsername(state,payload){

              state.loggedUsername=payload.username;
        },
        setUserKey(state,payload){
            state.loggedUserPrivateKey=payload.privateKey;
        },
        setConnection(state,payload){
            state.conn=payload.conn;
        },
    },
    actions:{
        addUsername(context,payload){
            context.commit('setUsername',payload)
        },
        addUserKey(context,payload){
            context.commit('setUserKey',payload)
        },
        addConnection(context,payload){
            context.commit('setConnection',payload)
        },
    },
    getters:{
        getUsername(state){
            return state.loggedUsername;
        },
        getPk(state){
            return state.loggedUserPrivateKey;
        },
        getConnection(state){
            return state.conn;
        },
    }
}