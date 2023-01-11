<script>
export default {
  data() {
    return {
      passcode: "",
      showWarning: false,
    };
  },
  methods: {
    passCheck() {
      this.showWarning = false;
      const queryString = window.location.search;
      const urlParams = new URLSearchParams(queryString);
      const redirectURI = urlParams.get("redirect_uri");

      const XHR = new XMLHttpRequest();
      const urlEncodedDataPairs = [];
      urlEncodedDataPairs.push(`password=${encodeURIComponent(this.passcode)}`);
      const urlEncodedData = urlEncodedDataPairs.join("&").replace(/%20/g, "+");
      XHR.addEventListener("load", (event) => {
        if (XHR.status == 200) {
          this.$router.push(decodeURIComponent(redirectURI));
        } else {
          this.showWarning = true;
        }
      });
      XHR.open("POST", "/api/login");
      XHR.setRequestHeader("Content-Type", "application/x-www-form-urlencoded");
      XHR.send(urlEncodedData);
    },
  },
};
</script>
<style>
#app {
  flex-direction: column;
}
.title {
  margin-top: 0;
}
.login-box {
  background: rgba(245, 245, 220, 0.638);
  border: black 1px solid;
  padding: 2rem;
  margin: 10px;
  border-radius: 2rem;
  align-items: center;
  text-align: center;
}
.formfield {
  display: block;
  width: 100%;
  margin-bottom: 1rem;
  border-radius: 5px;
  border: 1px solid gray;
}
input.formfield {
  height: 3rem;
  font-size: 1.5rem;
  box-sizing: border-box;
  text-align: center;
}
button.formfield {
  height: 3rem;
  font-size: 2rem;
}
</style>
<template>
  <div class="login-box">
    <h1 class="title">Enter Passcode</h1>
    <p v-if="showWarning">Incorrect password!</p>
    <input
      class="formfield"
      type="text"
      v-model="passcode"
      placeholder="Say the magic word!"
    />
    <button class="formfield" @click="passCheck">Log In</button>
  </div>
</template>
