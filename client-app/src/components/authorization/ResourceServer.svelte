<script lang="ts">
  import ScopeList from "./ScopeList.svelte";
  import type { ResourceServer } from "./models";

  let open = true;

  export let resourceServer: ResourceServer;
</script>

<div class="">
  <button
    id={resourceServer.id}
    type="button"
    aria-expanded={open}
    aria-controls={`${resourceServer.id}-scopes`}
    class="flex w-full items-center px-2 py-1 border border-stone-500 rounded-lg -outline-offset-4"
    on:click={() => open = !open}
  >
    <img class="w-6 h-6 mr-2 bg-white rounded-md inline-block" src={resourceServer.src} alt={`resource-server-${resourceServer.id}`}>
    <h3 class="inline-block text-lg">{resourceServer.title}</h3>
    <svg class={`ml-auto fill-stone-50 ${open ? '' : 'rotate-180'}`} xmlns="http://www.w3.org/2000/svg" height="24" viewBox="0 -960 960 960" width="24"><path d="m283-345-43-43 240-240 240 239-43 43-197-197-197 198Z"/></svg>
  </button>
  <div class={open ? '' : 'hidden'} id={`${resourceServer.id}-scopes`} aria-labelledby={resourceServer.id}>
    <ScopeList scopes={resourceServer.scopes} />
  </div>
</div>
