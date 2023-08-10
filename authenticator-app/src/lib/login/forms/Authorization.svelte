<script lang="ts">
  import ResourceServerList from "../authorization/ResourceServerList.svelte";
  import type { ResourceServer } from "../authorization/models";

  const resourceServers: ResourceServer[] = [
    {
      id: "resource-server-a",
      title: "Resource Server A",
      src: "/resource-server-a.png",
      scopes: [
        {
          id: 'read-profile',
          description: 'Read profile',
        },
        {
          id: 'write-profile',
          description: 'Write profile',
        },
      ]
    },
    {
      id: "resource-server-b",
      title: "Resource Server B",
      src: "/resource-server-b.png",
      scopes: [
        {
          id: 'cloud-storage',
          description: 'Access and modify cloud storage',
        },
      ]
    },
    {
      id: "resource-server-c",
      title: "Resource Server C",
      src: "/resource-server-c.png",
      scopes: [
        {
          id: 'test',
          description: 'Test application',
        },
        {
          id: 'build',
          description: 'Build application',
        },
        {
          id: 'deploy',
          description: 'Deploy application',
        },
      ]
    }
  ]

  export let submit: (scopes: string[]) => void;

  const onSubmit = (event: SubmitEvent) => {
    event.preventDefault();

    const formTarget = event.target as HTMLFormElement;

    const filteredScopes = resourceServers
      // map resource servers to their scopes
      .map((rs) => rs.scopes.map((s) => s.id))
      // flatten the array
      .flat()
      // filter out all scopes that are not checked
      .filter((s) => formTarget.elements[s].checked);
    
    submit(filteredScopes);
  }
</script>

<form class="flex flex-col flex-1 space-y-1" on:submit={onSubmit}>
  <h2 class="text-xl">
    Authorization
  </h2>
  <p>
    Client <em class="italic">"Modern Auth Demo"</em> is requesting access to your account. The following scopes are requested:
  </p>
  <ResourceServerList class="!my-2 max-h-96 overflow-y-scroll" resourceServers={resourceServers} />
  <button type="submit" class="self-end btn btn-primary">
    Authorize
  </button>
</form>