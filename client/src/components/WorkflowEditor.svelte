<script lang="ts">
    import ActionCard from "./ActionCard.svelte";
    import Connector from "./Connector.svelte";
    type Action = {
        text: string,
        logo: string
    }
    
    let actions: Action[] = [
        {
            text: "@darthtrooper goes online",
            logo: "i-logos-twitch"
        },
        {
            text: "Send email to list",
            logo: "i-logos-google-gmail"
        }
    ]
</script>

<section class="mt-32 w-full flex flex-col items-center">
    {#each actions as action, index}
        <ActionCard text={action.text} logo={action.logo}/>
        {#if index != actions.length - 1}
            <div class="w-[1px] h-8 bg-secodaryLight"></div>
            <Connector/>
            <div class="w-[1px] h-8 bg-secodaryLight"></div>
        {/if}
    {/each}
</section>
