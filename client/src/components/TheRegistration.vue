<template>
  <link rel="stylesheet" href="https://stackpath.bootstrapcdn.com/bootstrap/4.5.2/css/bootstrap.min.css"
        integrity="sha384-JcKb8q3iqJ61gNV9KGb8thSsNjpSL0n8PARn9HuZOnIxN0hoP+VmmDGMN5t9UJ0Z" crossorigin="anonymous">

  <section class="login ">
    <div class="login_box  ">
      <div class="main_box ">
        <div class="top_link"><a href="#"><img
            src="https://icon2.cleanpng.com/20180425/fue/kisspng-computer-icons-button-5ae031e2bf0904.7825510515246422747825.jpg"
            alt="">Return home</a></div>
        <div class="contact">
          <form action="" @submit.prevent="register">
            <h3>Registration</h3>
            <input type="text" placeholder="USERNAME" v-model="Username">
            <input type="text" placeholder="PASSWORD" v-model="Password">
            <button class="submit">Register</button>
          </form>
        </div>
      </div>

    </div>
  </section>


</template>

<script>


const SERVER_PUB_KEY = { x: "919a3e190688a262da21239706c80a4a5b11f897514bc9dcbae3aa627f72ade5", y: "365e1e475d29cafbc09e72e2147633e535f61c6a6b2edf456b08c7d2abec8290" };


export default {
  data() {
    return {
      Username: "",
      Password: "",
      clicked: false,
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
    register() {
      var ec = new this.$elliptic.ec('p256');
      var serv_key = ec.keyFromPublic(SERVER_PUB_KEY,'hex')
      if (this.Password===""){
        alert("Password can't be empty")
        throw new Error("Password can't be empty")

      }else if (this.Password.length<8){
        alert("Password should have length 8 minimum")
        throw new Error("Password should have length 8 minimum")
      }else if (this.Username===""){
        alert("Username can't be empty")
        throw new Error("Username can't be empty")
      }else if (this.Username.length<4 || this.Username.length>64){
        alert("Username should have length 4 minimum and can't exceed 64 symbols")
        throw new Error("Username should have length 4 minimum")
      }
      var regexPattern = /^[a-zA-Z0-9]+$/;
      if (!regexPattern.test(this.Username)){
        alert("Only english symbols and numbers are available");
        this.Username="";
        throw new Error("wrong symbols were used in username");
      }
      this.clicked = true;

      let randNonce = this.randN();
      fetch(
          "http://localhost:8081/api/preregister/user",
          {
            method: "POST",
            headers: {
              "Content-Type": "application/json",
            },
            body: JSON.stringify({
              Username: this.Username,
              UserNonce: randNonce,

            }),
          }
      )
          .then((response) => response.json()).then( data => {


              let hashed = this.$jscrypto.SHA256.hash(randNonce).toUint8Array();
           //   console.log(data)
              if (ec.verify(hashed,{r:data.r,s:data.s},serv_key)){

                var myPrivateKey=this.$jscrypto.PBKDF2.getKey(this.Password,data.salt,{keySize: 8,iterations: 150000})
                var keys = ec.keyFromPrivate(myPrivateKey.toString())

                let signedMessage = keys.sign(this.$jscrypto.SHA256.hash(data.nonce).toUint8Array())
              /*  console.log(signedMessage)*/
                //console.log(keys.getPublic().getX().toString())
                //console.log(keys.getPublic().getY().toString())
                //console.log(signedMessage.r.toString())
                //console.log(signedMessage.s.toString())
                fetch(
                    "http://localhost:8081/api/register/user",
                    {
                      method: "POST",
                      headers: {
                        "Content-Type": "application/json",
                      },
                      body: JSON.stringify({
                        Username: this.Username,
                        r:signedMessage.r.toString(),
                        s:signedMessage.s.toString(),
                        x:keys.getPublic().getX().toString(),
                        y:keys.getPublic().getY().toString(),
                      }),
                    }
                ).then((response) => {
                  if (!response.ok){
                    throw new Error("Something bad happened during registration");
                  }else{
                    this.$router.replace("/");
                  }
                })
              }
          })
          .catch((error) => {
            alert("Such user already exists!");
            console.log(error);
          }
          );/**/

    }
  }
}
</script>


<style scoped>

.login {
  height: 100vh;
  width: 100%;
  background: url("https://wallpapercrafter.com/sizes/1920x1080/225657-in-the-night-sky-full-of-stars-stands-a-bright-lon.jpg");
  position: relative;
}

.login_box {
  width: 600px;
  height: 600px;
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  background: none;
  border-radius: 10px;

  display: flex;
  overflow: hidden;
}

.login_box .main_box {
  width: 100%;
  height: 100%;
  padding: 25px 25px;

}

.main_box .top_link a {
  color: #452A5A;
  font-weight: 400;
}

.main_box .top_link {
  height: 20px
}

.main_box .contact {
  display: flex;
  align-items: center;
  justify-content: center;
  align-self: center;
  height: 100%;
  width: 73%;
  margin: auto;
}

.main_box h3 {
  text-align: center;
  margin-bottom: 40px;
}

.main_box input {
  border: none;

  margin: 15px 0;
  border-bottom: 1px solid #4f30677d;
  padding: 7px 9px;
  width: 100%;
  overflow: hidden;
  background: transparent;
  font-weight: 600;
  font-size: 14px;
}

.main_box {
  background: rgba(255, 255, 255, 0.95);
}

.submit {
  border: none;
  padding: 15px 70px;
  border-radius: 8px;
  display: block;
  margin: auto;
  margin-top: 120px;
  background: #583672;
  color: #fff;
  font-weight: bold;
  -webkit-box-shadow: 0px 9px 15px -11px rgba(88, 54, 114, 1);
  -moz-box-shadow: 0px 9px 15px -11px rgba(88, 54, 114, 1);
  box-shadow: 0px 9px 15px -11px rgba(88, 54, 114, 1);
}


.top_link img {
  width: 28px;
  padding-right: 7px;
  margin-top: -3px;
}


</style>