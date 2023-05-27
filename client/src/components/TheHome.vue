<template>
  <link
      href="https://cdn.jsdelivr.net/npm/bootstrap@5.0.2/dist/css/bootstrap.min.css"
      rel="stylesheet"
      integrity="sha384-EVSTQN3/azprG1Anm3QDgpJLIm9Nao0Yz1ztcQTwFspd3yD65VohhpuuCOmLASjC"
      crossorigin="anonymous"
  />


  <div class="backg">
    <div class="container">
      <div class="main">

        <form @submit.prevent="login">
          <label>
            <input class="inp" type="text" placeholder="Enter Your Username" v-model="Username">
          </label><br>
          <label>
            <input class="inp" type="text" placeholder="Enter Your Password" v-model="Password">
          </label>
          <br>
          <label>
            <button type="submit" class="btn btn-secondary">Log in</button></label>
            <br><br>
            <label><router-link class="x" to="/registration">
              <button class="btn btn-secondary">Register</button>
            </router-link></label>

        </form>

      </div>
    </div>
  </div>


</template>

<script>
const SERVER_PUB_KEY = { x: "dbd8137db1ae0ddeb60e4cad188963c9766de468eb2f1c8520c05b121c53bd5b", y: "4f31fa54cbed9e1cd5e7f2c9855e00ebb4db47cba7bcf5667363957278204747" };



export default {
 /* setup () {
    const store = useStore()

    function sendMessage () {
      const trimedText = text.value.trim()
      if (trimedText) {
        store.dispatch('sendMessage', {
          text: trimedText,
          thread: thread.value
        })
        this.text = ''
      }
    }

    return {
      sendMessage
    }
  },*/
  data() {
    return {
      Username: "",
      Password: "",
    }
  },
  methods: {
    randN() {
      let UserNonce = "";
      let nums = new Uint32Array(8);
      crypto.getRandomValues(nums);
      for (let i = 0; i < nums.length; i++) {
        UserNonce += nums[i].toString(16);
      }
      return UserNonce
    },
    login() {
      const sss=this;
      let nonce = this.randN();
      // alert()
      //this.clicked=true;
      fetch(
          "http://localhost:8081/api/login/user",
          {
            method: "POST",
            headers: {
              "Content-Type": "application/json",
            },
            body: JSON.stringify({
              Username: this.Username,
              UserNonce: nonce,
            }),
          }
      ).then((response) => response.json()).then(data => {
        var ec = new this.$elliptic.ec('p256');
        var serv_key = ec.keyFromPublic(SERVER_PUB_KEY,'hex')

        let hashed = this.$jscrypto.SHA256.hash(nonce).toUint8Array();
       // console.log(data)
        if (ec.verify(hashed, {r: data.r, s: data.s}, serv_key)) {

          var myPrivateKey = this.$jscrypto.PBKDF2.getKey(this.Password, data.salt, {keySize: 8, iterations: 150000})
          var keys = ec.keyFromPrivate(myPrivateKey.toString())
          //console.log("nonce    "+data.nonce)
          let signedMessage = keys.sign(this.$jscrypto.SHA256.hash(data.nonce).toUint8Array())
          let connection = new WebSocket("ws://localhost:8081/api/ws/"+this.Username+"/"+signedMessage.r.toString()+"_"+signedMessage.s.toString())

 //         connection.onmessage = function(event) {
 //           console.log(event);
 //         }

          let usrn = this.Username;
          connection.onopen = function() {
           // console.log("Successfull connection")
            sss.$store.dispatch('user/addConnection',{conn: connection})
            sss.$store.dispatch('user/addUserKey',{privateKey: keys})
            sss.$store.dispatch('user/addUsername', {username: usrn})
            sss.$router.push("/user")
          }
          connection.onclose = function(event){

            console.log(event)
            console.log("closed")
          }
        }
      })
          .catch((error) => {
            throw error;
              //  console.log(error);
              }
          );

      // console.log(this.loggedPersonMail)


    }
  }
}
</script>
<style scoped>
/* html, body {
    margin: 0;
    padding: 0;
  } */

.header {
  background-color: black;
}

.x {
  text-decoration: none;
  color: white;
}

.inp {
  padding: 10px 50px;
  border: 0;
  margin-bottom: 15px;
  border-radius: 10px;
}

.inp::placeholder {
  font-family: Arial, Helvetica, sans-serif;
  text-align: center;


}

.container {
  max-width: 940px;
  margin: 0 auto;
  padding: 0 10px;
  /* padding: 0 ; */
}

.backg {
  background: url(https://images.wallpaperscraft.com/image/single/stars_starry_sky_night_182857_1920x1080.jpg);
  background-size: cover;
  background-position: center center;
  background-repeat: no-repeat;
  height: 100vh;
}

.nav {
  margin: 0;
  padding: 20px 0;

}

.nav li {
  display: inline;
  color: rgb(19, 18, 18);
  font-family: 'Raleway', sans-serif;
  font-weight: 600;
  font-size: 12px;
  text-transform: uppercase;
  margin-left: 10px;
  margin-right: 10px;
}

.main {
  position: relative;
  top: 280px;
  text-align: center;
}


.btn-main {
  background-color: #333;
  color: #fff;
  font-family: 'Raleway', sans-serif;
  font-weight: 600;
  font-size: 18px;
  letter-spacing: 1.3px;
  padding: 16px 40px;
  text-decoration: none;
  text-transform: uppercase;
}

.btn-default {
  color: #333;
  border: 1px solid #333333;
  font-family: 'Raleway', sans-serif;
  font-weight: 600;
  font-size: 10px;
  letter-spacing: 1.3px;
  padding: 10px 20px;
  text-decoration: none;
  text-transform: uppercase;
  display: inline-block;
  margin-bottom: 20px;
  padding-right: 50px;
  padding-left: 50px;
}

.supporting {
  padding-top: 80px;
  padding-bottom: 100px;
}

.supporting .col {
  float: left;
  width: 33%;
  font-family: 'Raleway', sans-serif;
  text-align: center;
}

.supporting img {
  height: 32px;
}

.supporting h2 {
  font-weight: 600;
  font-size: 23px;
  text-transform: uppercase;
}

.supporting p {
  font-weight: 400;
  font-size: 14px;
  line-height: 20px;
  padding: 0 50px;
  margin-bottom: 40px;
}

.clearfix {
  clear: both;
}

.footer {
  background-color: #333;
  color: #fff;
  padding: 10px 0;
}

.footer p {
  font-family: 'Raleway', sans-serif;
  text-transform: uppercase;
  font-size: 15px;
  margin-top: 15px;
  text-align: center;
}

@media (max-width: 500px) {
  .main h1 {
    font-size: 50px;
    padding: 0 40px;
  }

  .supporting .col {
    width: 100%;
  }
}
</style>


