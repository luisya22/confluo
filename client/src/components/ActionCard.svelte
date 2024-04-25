<script lang="ts">
    export let text: string;
    export let logo: string;
</script>

<div class="px-10 py-10 w-1/2 border border-cardBorder rounded-xl">
    <div class="flex justify-start items-center">
        <div class="w-1/12">
            <span class={`w-10 h-10 ${logo}`}></span>
        </div>
        <div class="w-11/12">
            <h4 class="text-secodaryLight text-2xl font-bold px-10 w-11/12">{text}</h4>
        </div>
    </div>
</div>
