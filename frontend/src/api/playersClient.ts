import { create } from "@bufbuild/protobuf"
import { createClient } from "@connectrpc/connect"
import { createConnectTransport } from "@connectrpc/connect-web"
import { applyDevAdminHeaders } from "@/api/devAdmin"
import {
  CreatePlayerRequestSchema,
  DeleteAllPlayersRequestSchema,
  DeletePlayerRequestSchema,
  GenerateRandomPlayersRequestSchema,
  ListPlayersRequestSchema,
  PlayerService,
  UpdatePlayerRequestSchema,
} from "@/gen/cribbly/v1/players_pb"

const transport = createConnectTransport({
  baseUrl: "/api",
  fetch: (input, init) => {
    const headers = new Headers(init?.headers)
    applyDevAdminHeaders(headers)
    return fetch(input, { ...init, headers, credentials: "include" })
  },
})

const client = createClient(PlayerService, transport)

export async function listPlayers() {
  return client.listPlayers(create(ListPlayersRequestSchema, {}))
}

export async function createPlayer(firstName: string, lastName: string) {
  return client.createPlayer(
    create(CreatePlayerRequestSchema, { firstName, lastName }),
  )
}

export async function updatePlayer(
  id: string,
  firstName: string,
  lastName: string,
) {
  return client.updatePlayer(
    create(UpdatePlayerRequestSchema, { id, firstName, lastName }),
  )
}

export async function deletePlayer(id: string) {
  return client.deletePlayer(create(DeletePlayerRequestSchema, { id }))
}

export async function deleteAllPlayers() {
  return client.deleteAllPlayers(create(DeleteAllPlayersRequestSchema, {}))
}

export async function generateRandomPlayers(count: number) {
  return client.generateRandomPlayers(
    create(GenerateRandomPlayersRequestSchema, { count }),
  )
}
