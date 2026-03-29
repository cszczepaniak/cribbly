import { Link, type LinkProps } from "react-router-dom"
import { withReactQuery } from "../reactNav"

type Props = Omit<LinkProps, "to"> & { to: string }

/** Same as {@link Link}, but always includes `?react=true` for the Go server. */
export function ReactLink({ to, ...rest }: Props) {
  return <Link to={withReactQuery(to)} {...rest} />
}
