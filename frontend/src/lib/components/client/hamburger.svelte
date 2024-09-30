<script lang="ts">
    export let activated: boolean;
    export let toggle: () => void;
</script>

<button class="relative group" on:click={toggle}>
    <div class="relative flex overflow-hidden items-center justify-center rounded-full w-[50px] h-[50px] transform transition-all ring-opacity-30 duration-200">
        <div class="flex flex-col justify-between w-[20px] h-[20px] transform transition-all duration-300 origin-center overflow-hidden {activated ? '-translate-x-1.5 rotate-180' : ''}">
            <div class="bg-white h-[2px] w-7 transform transition-all duration-300 origin-left delay-150 {activated ? 'rotate-[42deg] w-[70%]' : ''}"></div>
            <div class="bg-white h-[2px] w-7 rounded transform transition-all duration-300 {activated ? 'translate-x-10' : ''}"></div>
            <div class="bg-white h-[2px] w-7 transform transition-all duration-300 origin-left delay-150 {activated ? '-rotate-[42deg] w-[70%]' : ''}"></div>
        </div>
    </div>
</button>