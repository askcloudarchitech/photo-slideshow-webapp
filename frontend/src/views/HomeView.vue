<script>
import axios from "axios";
export default {
  data() {
    return {
      appearance: "default",
      showWarning: false,
      selectedFile: null,
      warningMessage: "There was an error. Try again",
    };
  },
  methods: {
    onFileChanged(event) {
      this.selectedFile = event.target.files[0];
      this.appearance = "imgSelected";
      this.showWarning = false;
    },
    async onUpload(event) {
      this.appearance = "uploading";
      this.showWarning = false;
      const formData = new FormData();
      formData.append("file", this.selectedFile, this.selectedFile.name);
      try {
        let res = await axios.post("/api/upload", formData);
        this.warningMessage = "Success!";
        this.showWarning = true;
        this.appearance = "default";
      } catch (err) {
        this.appearance = "default";
        this.warningMessage =
          "There was an error. Try again " + err.response.status;
        this.showWarning = true;
      }
    },
    onCancel(event) {
      this.selectedFile = null;
      this.showWarning = false;
      this.appearance = "default";
    },
  },
};
</script>

<style>
.upload-box {
  background: rgba(245, 245, 220, 0.638);
  border: black 1px solid;
  padding: 2rem;
  margin: 10px;
  border-radius: 2rem;
  align-items: center;
  text-align: center;
}
input[type="file"] {
  display: none;
}
.custom-file-upload {
  border: 1px solid #ccc;
  display: inline-block;
  padding: 15px 30px;
  cursor: pointer;
  background-color: lightgray;
  border-radius: 15px;
  font-size: 2rem;
}
#selected-image {
  padding-bottom: 20px;
}
</style>

<template>
  <div class="upload-box">
    <h1 class="title">Add an Image</h1>
    <p v-if="showWarning">{{ warningMessage }}</p>
    <div v-if="appearance === 'imgSelected'" id="selected-image">
      Image: {{ selectedFile.name }}
    </div>
    <div v-if="appearance === 'uploading'">Please Wait...</div>
    <label v-if="appearance === 'default'" class="custom-file-upload">
      <input
        v-if="appearance === 'default'"
        class="formfield"
        type="file"
        @change="onFileChanged"
      />
      Select Image
    </label>
    <button
      v-if="appearance === 'imgSelected'"
      class="formfield"
      @click="onUpload"
    >
      Send
    </button>
    <button
      v-if="appearance === 'imgSelected'"
      class="formfield"
      @click="onCancel"
    >
      Cancel
    </button>
  </div>
</template>
