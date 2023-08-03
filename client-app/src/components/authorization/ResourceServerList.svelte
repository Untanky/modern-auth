<script lang="ts">
  import ResourceServer from "./ResourceServer.svelte";
  import type { ResourceServer as ResourceServerType } from "./models";

  export let resourceServers: ResourceServerType[];
  let klass: string = "";
  export { klass as class };
</script>

<ul class="space-y-2.5 {klass}">
  {#each resourceServers as resourceServer}
  <li>
    <ResourceServer resourceServer={resourceServer}></ResourceServer>
  </li>
  {/each}
</ul>