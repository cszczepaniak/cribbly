import { create } from '@bufbuild/protobuf'
import { createClient } from '@connectrpc/connect'
import { createConnectTransport } from '@connectrpc/connect-web'
import { RoomCodeService, SetRoomCodeRequestSchema } from '@/gen/cribbly/v1/roomcode_pb'

const transport = createConnectTransport({
  baseUrl: '/api',
  fetch: (input, init) => fetch(input, { ...init, credentials: 'include' }),
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
