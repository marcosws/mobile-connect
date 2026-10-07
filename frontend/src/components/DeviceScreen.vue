<script setup>
import { ref, onMounted, onUnmounted } from "vue";

let controlSocket = null;

let startX = 0
let startY = 0

let currentX = 0
let currentY = 0

let isDragging = false

let pressStart = 0

const threshold = 20
const LONG_PRESS_TIME = 700

const props = defineProps({
  deviceId: {
    type: String,
    required: true,
  },
});

const imageUrl = ref("");

let videoSocket = null;
let currentUrl = null;


// Flag para sabermos que o componente está sendo desmontado
// e qualquer evento posterior (onerror/onclose) deve ser silencioso.
let isUnmounting = false;


onMounted(() => {

    isUnmounting = false;

  videoSocket = new WebSocket(
    `ws://localhost:8080/devices/${props.deviceId}/stream/video`
  );

  videoSocket.binaryType = "arraybuffer";

  videoSocket.onopen = () => {
    console.log("Video connected");
  };

  videoSocket.onmessage = (event) => {

    if (currentUrl) {
      URL.revokeObjectURL(currentUrl);
    }

    const blob = new Blob(
      [event.data],
      {
        type: "image/png",
      }
    );

    currentUrl =
      URL.createObjectURL(blob);

    imageUrl.value =
      currentUrl;
  };

  videoSocket.onerror = (err) => {
    // Se já estamos desmontando ou fechando, ignora o erro:
    // é o navegador abortando a conexão intencionalmente.
    if (isUnmounting) return;
    console.error("Erro no videoSocket:", err);

  };
  videoSocket.onclose = (event) => {
    if (!event.wasClean && !isUnmounting) {
      console.warn(
        `videoSocket fechado inesperadamente: code=${event.code}`
      );
    }
  };

    controlSocket = new WebSocket(
    `ws://localhost:8080/devices/${props.deviceId}/stream/control`
    );

    controlSocket.onopen = () => {
    console.log("Control connected");
    };

    controlSocket.onerror = (err) => {
    if (isUnmounting) return;
    console.error("Erro no controlSocket:", err);
  };
 
  controlSocket.onclose = (event) => {
    if (!event.wasClean && !isUnmounting) {
      console.warn(
        `controlSocket fechado inesperadamente: code=${event.code}`
      );
    }
  };


});

onUnmounted(() => {

    isUnmounting = true;

  if (currentUrl) {
    URL.revokeObjectURL(currentUrl);
  }

  closeSocketSafely(controlSocket);
  closeSocketSafely(videoSocket);

  if (currentUrl) URL.revokeObjectURL(currentUrl);
 
  controlSocket = null;
  videoSocket = null;


//   controlSocket?.close();
//   videoSocket?.close();

});




// --- Helper de fechamento seguro -------------------------------------
function closeSocketSafely(socket) {
  if (!socket) return;
 
  switch (socket.readyState) {
    case WebSocket.OPEN:
      socket.close(1000, "Componente desmontado");
      break;
 
    case WebSocket.CONNECTING:
      // Ainda conectando: não dá pra fechar "limpo" agora.
      // Esperamos abrir e fechamos imediatamente em seguida.
      socket.addEventListener(
        "open",
        () => socket.close(1000, "Cancelado antes de usar"),
        { once: true }
      );
      break;
 
    // CLOSING ou CLOSED: nada a fazer
    default:
      break;
  }
}




// function sendTap(event) {

//   const rect = event.target.getBoundingClientRect()

//   const scaleX =
//     event.target.naturalWidth / rect.width

//   const scaleY =
//     event.target.naturalHeight / rect.height

//   const x = Math.round(
//     (event.clientX - rect.left) * scaleX
//   )

//   const y = Math.round(
//     (event.clientY - rect.top) * scaleY
//   )

//   console.log("TAP:", { x, y })

//   controlSocket.send(
//     JSON.stringify({
//       type: "tap",
//       x,
//       y
//     })
//   )
// }
// Swipe handling =========================================================
function endSwipe(event) {

  if (!isDragging) return

  isDragging = false

  const pressDuration = Date.now() - pressStart

  console.log({
    pressDuration
  })

  window.removeEventListener("mousemove", moveSwipe)
  window.removeEventListener("mouseup", endSwipe)

  const dx = currentX - startX
  const dy = currentY - startY

  const distance = Math.sqrt(dx * dx + dy * dy)

  console.log("SWIPE RESULT:", {
    startX,
    startY,
    currentX,
    currentY,
    distance
  })

  if (distance < threshold) {

    if (pressDuration >= LONG_PRESS_TIME) {

      console.log("LONG PRESS")

      controlSocket.send(
        JSON.stringify({
          type: "longpress",
          x: startX,
          y: startY,
          duration: pressDuration
        })
      )

      console.log("After LONG PRESS")

      return
    }

    console.log("TAP")

    controlSocket.send(
      JSON.stringify({
        type: "tap",
        x: startX,
        y: startY
      })
    )

    return
  }

  controlSocket.send(JSON.stringify({
    type: "swipe",
    x: startX,
    y: startY,
    x2: currentX,
    y2: currentY,
    duration: 300
  }))
}

function moveSwipe(event) {

  if (!isDragging) return

  const img = document.querySelector(".device-screen")
  if (!img) return

  const rect = img.getBoundingClientRect()

  const scaleX = img.naturalWidth / rect.width
  const scaleY = img.naturalHeight / rect.height

  currentX = Math.round((event.clientX - rect.left) * scaleX)
  currentY = Math.round((event.clientY - rect.top) * scaleY)
}

function startSwipe(event) {

  pressStart = Date.now()

  event.preventDefault()

  const rect = event.target.getBoundingClientRect()

  const scaleX = event.target.naturalWidth / rect.width
  const scaleY = event.target.naturalHeight / rect.height

  startX = Math.round((event.clientX - rect.left) * scaleX)
  startY = Math.round((event.clientY - rect.top) * scaleY)

  currentX = startX
  currentY = startY

  isDragging = true

  window.addEventListener("mousemove", moveSwipe)
  window.addEventListener("mouseup", endSwipe)
}

</script>

<template>

  <div class="remote-panel">

    <div class="screen-container">

      <img
      v-if="imageUrl"
      :src="imageUrl"
      class="device-screen"
      draggable="false"
      @dragstart.prevent
      @pointerdown="startSwipe"
      @pointermove="moveSwipe"
      @pointerup="endSwipe"
      />

    </div>

  </div>

</template>

<style scoped>
.screen-container {
  display: flex;
  justify-content: center;
  border: 2px solid #ccc;
  border-radius: 8px;
  padding: 4px;
  margin-left: 10px;
  margin-right: 10px;
}

/* .device-screen {

    max-height: 800px;
  border: 1px solid #ccc;
  -webkit-user-drag: none;
  user-select: none;
  touch-action: none;
} */

.device-screen {
  width: auto;      /* ou 360px se quiser estilo mobile */
  height: auto;
  display: block;
  max-height: 70vh;

  object-fit: contain;

  -webkit-user-drag: none;
  user-select: none;
  touch-action: none;
  border: 2px solid #ccc;
  border-radius: 8px;
}
</style>