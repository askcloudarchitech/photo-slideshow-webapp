<template>
  <div class="date-container" :data-vis="visible">{{ relativeTime }}</div>
  <div class="img-container" :data-vis="visible">
    <img class="photo" :src="image" />
  </div>
</template>

<script>
export default {
  props: {
    image: String,
    visible: Boolean,
    relativeTime: String,
  },
};
</script>

<style>
.photo {
  max-height: 100%;
  width: auto;
  margin: auto;
  display: block;
}
.img-container {
  position: absolute;
  bottom: 0;
  width: 100%;
  height: 100%;
}
.date-container {
  position: absolute;
  top: 60px;
  left: 100px;
  background-color: #000;
  padding: 1rem;
  border: 1px solid;
  border-radius: 100px;
  z-index: 10000;
  color: wheat;
  font-size: 4rem;
}
[data-vis="true"] {
  transition: opacity 2s;
  opacity: 1;
}
[data-vis="false"] {
  transition: opacity 2s;
  opacity: 0;
}
</style>
