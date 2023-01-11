<script>
import SlideshowImage from "../components/SlideshowImage.vue";
import RelativeTime from "@yaireo/relative-time";
const relativeTime = new RelativeTime();
export default {
  components: {
    SlideshowImage,
  },
  data() {
    return {
      images: [
        {
          image: "/photos/none.jpg",
          vis: false,
          relativeTime: "yesterday",
        },
        {
          image: "/photos/none2.jpg",
          vis: true,
          relativeTime: "yesterday",
        },
      ],
      visibleImage: 0,
    };
  },
  mounted() {
    setInterval(async () => {
      let hiddenimage = this.visibleImage === 0 ? 1 : 0;
      let res = await fetch("/api/slideshow/next");
      let data = await res.json();
      this.images[hiddenimage].image = data.Name;
      this.images[hiddenimage].relativeTime = relativeTime.from(
        new Date(data.TimeTaken * 1000)
      );
      setTimeout(() => {
        let hiddenimage = this.visibleImage === 0 ? 1 : 0;
        this.images[hiddenimage].vis = true;
        this.images[this.visibleImage].vis = false;
        this.visibleImage = hiddenimage;
      }, 2000);
    }, 7000);
  },
};
</script>

<template>
  <div class="frame">
    <div class="imageholder">
      <SlideshowImage
        :image="this.images[0].image"
        :visible="this.images[0].vis"
        :relativeTime="this.images[0].relativeTime"
      />

      <SlideshowImage
        :image="this.images[1].image"
        :visible="this.images[1].vis"
        :relativeTime="this.images[1].relativeTime"
      />
    </div>
  </div>
</template>

<style>
body {
  height: 100vh;
  width: 100vw;
  margin: 0px;
  background: repeating-linear-gradient(90deg, #cc231e, #0f8a5f, #cc231e);
  background-repeat: x;
  background-size: 25% 100%;
  animation: gradient 30s linear infinite;
}

@keyframes gradient {
  0% {
    background-position: 100% 50%;
  }
  100% {
    background-position: -33% 50%;
  }
}

.imageholder {
  position: relative;
  width: 100%;
  height: 100%;
}

.frame {
  background-color: rgba(0, 0, 0, 0);
  height: 95%;
  width: 100%;
}
</style>
