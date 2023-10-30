<script lang="ts">
    export const callback: (card: string) => void = () => {};

    let socket = new WebSocket('ws://localhost:3737/');

    type WsData = {
        uid: string;
    };

    socket.onmessage = (event) => {
        const data: WsData = JSON.parse(event.data);
        callback(data.uid);
    };
</script>