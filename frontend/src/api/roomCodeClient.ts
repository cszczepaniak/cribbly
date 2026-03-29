import { create } from "@bufbuild/protobuf"
import { createClient, type CallOptions } from "@connectrpc/connect"
import { createConnectTransport } from "@connectrpc/connect-web"
import {
  CheckRoomAccessRequestSchema,
  RoomCodeService,
  SetRoomCodeRequestSchema,
} from "@/gen/cribbly/v1/roomcode_pb"

const transport = createConnectTransport({
  baseUrl: "/api",
  fetch: (input, init) => fetch(input, { ...init, credentials: "include" }),
})

const client = createClient(RoomCodeService, transport)

/**
 * Validates the room code via Connect; on success the server sets the HttpOnly `room_code`
 * cookie (same as legacy POST /room-code). Include credentials on all `/api` calls so the
 * cookie is sent on later navigations and RPCs.
 */
export async function setRoomCode(code: string) {
  return client.setRoomCode(create(SetRoomCodeRequestSchema, { code }))
}

/** Returns whether the current browser has a valid room cookie or admin session (same rules as the Go server). */
export async function checkRoomAccess() {
  return client.checkRoomAccess(create(CheckRoomAccessRequestSchema, {}))
}

export async function doSomething(options?: CallOptions) {
  for await (const res of client.doSomething({}, options)) {
    console.log(res.data)
  }
}
