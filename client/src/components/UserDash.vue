<template>
  <link
    href="https://cdn.jsdelivr.net/npm/bootstrap@5.0.2/dist/css/bootstrap.min.css"
    rel="stylesheet"
    integrity="sha384-EVSTQN3/azprG1Anm3QDgpJLIm9Nao0Yz1ztcQTwFspd3yD65VohhpuuCOmLASjC"
    crossorigin="anonymous"
  />

  <div class="row  g-0 z bg-danger">
    <nav class="navbar navbar-light bg-light hdr">
      <form class="form-inline">
        <input
            type="search"
            class="form-control srch"
            placeholder="User to chat with"
            aria-label="User to chat with"
            aria-describedby="basic-addon2"
            v-model="userToFind"
        />

      </form>
      <button @click="findUser()" class="btn btn-outline-success my-2 my-sm-0 srchbtn">Chat now!</button>
      <b class="chtrname">{{ chatterName }}</b>
      <button @click="logout" type="button" class="btn btn-danger lgtbtn">Logout</button>

    </nav>
    <div class="col-3 chat-user shadow-lg ">
        <chat-user
          v-for="chatter in chatters"
          :key="chatter.username"

          :username="chatter.username"
          :pubKey="chatter.pubKey"
          @startChat="startChat"
        >
       </chat-user>

      <div class="chtud">
        <button @click="loadUsers" v-if="!lastLoadedUser" type="button" class="btn btn-primary loadU">Load More Users</button>
      </div>
    </div>

    <div class="col-9 convo">
      <div class="col-12 all-msg ">

        <div v-if="chatterName && !lastLoadedMessage" class="chtrb">
          <button @click="loadMessages" type="button" class="btn btn-primary loadM">Load More Messages</button>
        </div>
       <div   v-if="chatterName">
        <the-chat
          v-for="(x, i) in chatterMessages"
          :key="i.toString()+x.username"
          :message="x.message"
          :username="x.username"
        ></the-chat>
       </div>
      </div>

      <form @submit.prevent="sendMessage" class="w-100 mb-5" action="">
        <!-- <div class=" child"> -->
        <div class="send-msg">

          <div>
            <button type="submit" class="btn btn-primary send-btn" v-if="chatterName">Send</button>
            <input v-if="chatterName"
              type="text"
              class="form-control inpt"
              placeholder="Enter Your Message"
              aria-label="Enter Your Message"
              aria-describedby="basic-addon2"
              v-model="Text"
            />
          </div>
          <div>

          </div>
        </div>
        <!-- </div> -->
      </form>
      <!-- </div> -->

      <!-- <form action="">
          <input type="text" class="inpt" placeholder="Enter Your Message">
          <button type="submit" class="btn btn-primary p-2 ">Send</button>
        </form> -->
    </div>
  </div>
</template>

<script>

import ChatUser from "./ChatUser.vue";
import TheChat from "./TheChat.vue";
export default {
  components: {
    ChatUser,
    TheChat,
  },
  data() {
    return {
      userToFind: "",
      results: [], // users
      chatters: [],
      chatterMessages: [],
      results2: [], //messages assume ? used in conversation
      loggedUsername: this.$store.getters["user/getUsername"], // my email
      loggedUserKey: this.$store.getters["user/getPk"],
      Text: "", // my text
      chatterName: "", //chatter username
      chatterAESKey: "",
      conn: this.$store.getters["user/getConnection"],
      lastLoadedUser:false,
      lastLoadedMessage:false,
    };
  },
  mounted() {
    const sss = this;
   // console.log("mounted() username   ",sss.$store.getters["user/getUsername"]);
    sss.setWs()
  //  console.log("sss.setws finished in mounted ")
 //   sss.conn.onmessage = function(event){
      // to do
 //   }
  },
  methods: {
/*
 packet:
  type
  from
  to
  mess
* */
    setWs(){
      const sss=this;

      sss.conn.onmessage =  async function (event) {

        let data = JSON.parse(event.data)
      //  console.log(data)
        await sss.doSomethingWithData(data)
      }
      sss.loadUsers();
    },
    async doSomethingWithData(data) {
      const sss = this;
      if (data.type === 0) {
        if (sss.chatterName !== data.from && sss.chatterName !== "") {
          //console.log("we in first if")
        } else if (sss.chatterName === "") {
          //console.log("we in else if")
          await sss.setChatterData(data.from).then(dats=>{sss.chatterAESKey=dats});
          let decryptedMsg = sss.$cryptojs.AES.decrypt(data.message, sss.chatterAESKey).toString(sss.$cryptojs.enc.Utf8)
          sss.chatterMessages.push({
            message: decryptedMsg,
            username: data.from,
          })
        } else {
         // console.log("we in else ")
          let aes = sss.$cryptojs.AES
          let decryptedMsg = aes.decrypt(data.message, sss.chatterAESKey).toString(sss.$cryptojs.enc.Utf8)
          sss.chatterMessages.push({
            message: decryptedMsg,
            username: data.from,
          })
        }

      }
      else if (data.type=== 101){
        let users = JSON.parse(data.message);
    //    console.log("users     ",users)
        let lenUsers=users.length;
        for (let i = 0; i < lenUsers; i++) {
          sss.chatters.push({username: users[i].username,
            pubKey: users[i].key,})
        }
      }
      else if (data.type=== 102){
        sss.lastLoadedUser=true;
        let users = JSON.parse(data.message);
     //   console.log("users     ",users)
        let lenUsers=users.length;
        for (let i = 0; i < lenUsers; i++) {
          sss.chatters.push({username: users[i].username,
            pubKey: users[i].key,})
        }
      }
      else if (data.type=== 111){
        let msgs = JSON.parse(data.message);
        let lenUsers=msgs.length;
       // console.log("imagine ",sss.chatterAESKey)
        for (let i = 0; i < lenUsers; i++) {
      //   console.log(msgs[i].username)
          let decryptedMsg = sss.$cryptojs.AES.decrypt(msgs[i].mess, sss.chatterAESKey).toString(sss.$cryptojs.enc.Utf8)
          sss.chatterMessages.unshift({message: decryptedMsg,
            username: msgs[i].username,
            })
        }
      }
      else if (data.type=== 112){
        sss.lastLoadedMessage=true;
        let msgs = JSON.parse(data.message);
        let lenUsers=msgs.length;
     //   console.log("imagine ",sss.chatterAESKey)
        for (let i = 0; i < lenUsers; i++) {
     //     console.log(msgs[i].username)
          let decryptedMsg = sss.$cryptojs.AES.decrypt(msgs[i].mess, sss.chatterAESKey).toString(sss.$cryptojs.enc.Utf8)
          sss.chatterMessages.unshift({message: decryptedMsg,
            username: msgs[i].username,
          })
        }
      }
    //  console.log("message appeared")
    //  console.log(sss.chatterName)
    //  console.log(sss.chatterAESKey)

    },
    //load 5 last users on login
    loadUsers(){
      const sss = this;
      let pkt = {
        type: 100,
        from: sss.loggedUsername,
        to: "",
        message: "",
      }
      sss.conn.send(JSON.stringify(pkt))
     // console.log("sent load users")
    },
    loadMessages(){
      const sss = this;
      if (sss.chatterName===""){
        return;
      }
      let pkt = {
        type: 110,
        from: sss.loggedUsername,
        to: sss.chatterName,
        message: sss.chatterMessages.length.toString(),
      }
      sss.conn.send(JSON.stringify(pkt))
    },
    async setChatterData(name){
      return await this.getUser(name);
      /*console.log("done setting ", setKey)
      console.log("done setting ",this.chatterAESKey, this.chatterName)
      return setKey*/
    },
    async getUser(name){
      const sss = this;
      return fetch("http://localhost:8081/api/get/user/"+name, {
        method: "GET",
      }).then((response) => {
            if (!response.ok) {
             // console.log("IN GETUSER   1")
              alert("User wasn't get :(");
            }else {
            //  console.log("IN GETUSER   2")
              return response.json();
              //console.log(response)
            }
          }).then((data) => {

          //  console.log("IN GETUSER   3")
        var ec = new sss.$elliptic.ec('p256');

        let keysplit = data.key.split("_");
        sss.loggedUserKey.getPublic()
        let tempChatterKey = ec.keyFromPublic({x:keysplit[0],y:keysplit[1]},'hex');


        var shared = sss.loggedUserKey.derive(tempChatterKey.getPublic());
        let usr = {
          username: data.username,
          pubKey: data.key,
        };
        sss.chatters.push(usr)
        sss.chatterName = data.username;
        let sharedKey = sss.$jscrypto.SHA256.hash(shared.toString('hex'));
        let sharedKeyStr = sharedKey.toString()
        sss.chatterAESKey = sharedKeyStr;
       // console.log("IN GETUSER   ",sharedKeyStr)
        return sharedKeyStr
      })
          .catch((error) => {
            throw error;
           // console.log("get user :   ",error);
          })
    },
    sendMessage(){
      const sss = this;
      if (sss.Text===""){
        console.log("yes, return");
        return;
      }
      let aes = sss.$cryptojs.AES;
      let encrypted = aes.encrypt(sss.Text,sss.chatterAESKey);
   /*   console.log("key used to ecnrypt ",sss.chatterAESKey)
      console.log(encrypted)
      console.log(encrypted.toString())
      console.log("cryptojs tostr ", sss.$cryptojs.AES.encrypt(sss.Text,sss.chatterAESKey).toString())
      console.log("type of ", typeof(encrypted.toString()))*/
      let pkt = {
        type: 0,
        from: sss.loggedUsername,
        to: sss.chatterName,
        message: encrypted.toString(),
      }
      sss.conn.send(JSON.stringify(pkt))
      sss.chatterMessages.push({
        message:sss.Text,
        username:sss.loggedUsername,
      })
      sss.Text=""
    },
    findUser() {
      const sss = this;
      for (let i in sss.chatters){
        if (sss.chatters[i].username===sss.userToFind){
          alert("already found one")
          return
        }
      }
      if (sss.userToFind===sss.loggedUsername){
        alert("Don't speak to yourself...")
        return
      }
      if (sss.userToFind!==""){
        fetch("http://localhost:8081/api/get/user/"+sss.userToFind, {
          method: "GET",
        })
            .then((response) => {
              if (!response.ok) {
                alert("User wasn't found :(");
              }else {
                return response.json();
                //console.log(response)
              }
            }).then((data) => {
          //console.log(data);

          var ec = new sss.$elliptic.ec('p256');

          let keysplit = data.key.split("_");
          sss.loggedUserKey.getPublic()
          let tempChatterKey = ec.keyFromPublic({x:keysplit[0],y:keysplit[1]},'hex');


          var shared = sss.loggedUserKey.derive(tempChatterKey.getPublic());
          let usr = {
            username: data.username,
            pubKey: data.key,
          };
          sss.chatters.push(usr)
          sss.chatterName = data.username;
         // console.log("sh.tostr(hex)",shared.toString('hex'))
          let sharedKey = sss.$jscrypto.SHA256.hash(shared.toString('hex'));
         // console.log("shk ",sharedKey)
          let sharedKeyStr = sharedKey.toString()
        //  console.log("shkstr ",sharedKeyStr)
          sss.chatterAESKey = sharedKeyStr;
        })
            .catch((error) => {
              console.log("find user :   ",error);
            });
      }
    },
    logout(){
      const sss = this;
      sss.$router.replace('/');
      sss.results2=[];
      sss.conn.close();
      sss.$store.dispatch('user/addConnection',{conn: null});
      sss.$store.dispatch('user/addUserKey',{privateKey: null});
      sss.$store.dispatch('user/addUsername',{username: ""});
    },
    startChat: function (name, pk) {
      const sss = this;
      // alert();
      // console.log(value);
      var ec = new sss.$elliptic.ec('p256');

      let splitPk = pk.split("_");
      var chatterPk = ec.keyFromPublic({x: splitPk[0],y:splitPk[1]}, 'hex');
      var shared = sss.loggedUserKey.derive(chatterPk.getPublic())
      sss.chatterName = name;
      sss.chatterAESKey = sss.$jscrypto.SHA256.hash(shared.toString('hex')).toString();
      sss.loadMessages()
      // sss.Id2=value.Id;
      // console.log(sss.Id2)
      //console.log(sss.results2);
    },
  },

};
</script>

<style scoped>
.convo{
  height: 100vh;
}
.chtrname{
  margin-left: 10%;
  margin-right:5%;
}
.chtrb{
  text-align: center;
}
.loadM{
  margin-top:10px;
}
.loadU{
  margin-top:10px;
}
.hdr{
  alignment: left;
  justify-content: left;
}
.hdr .lgtbtn{
  margin-left: auto;
  margin-right:10px;
}
.srchbtn{
  margin-left: 10px;
}
.srch{
  margin-left: 5px;
}
.all-msg {
  height: 75%;
  overflow-y:auto;
}
.send-btn {
  /* margin-left: 1000px; */
  margin-right: 10px;
  float: right;
}
.cht {
  margin-top:5px;
}
.inpt {
  /* padding: 10px 350px; */
  border: 0;
  margin-bottom: 15px;
  margin-top:5%;
  margin-left: 20px;
  margin-right: 10px;
  width:90%;
}
.inpt::placeholder {
  font-family: Arial, Helvetica, sans-serif;
}
.chat-user {
  background: rgba(255, 255, 255, 0.562);
  overflow-y:auto;
  height: 100%;
}
.chtud{
  text-align: center;
}
.send-msg{
  overflow: hidden;
  /* background: yellow; */
}
img {
  height: 30px;
  width: 30px;
  border-radius: 50%;
}
.z {
  background: linear-gradient(#332042, #653d84);
  background-repeat: no-repeat;
  background-size: cover;
  background-attachment: fixed;
  background-position: center center;
  height: 100vh;
}
</style>
