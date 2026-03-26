import { useNavigate, type NavigateOptions } from 'react-router-dom'
import { withReactQuery } from '../reactNav'

/** Like {@link useNavigate}, but paths include `?react=true` for the Go server. */
export function useReactNavigate() {
  const navigate = useNavigate()
  return (to: string, opts?: NavigateOptions) =>
    navigate(withReactQuery(to), opts)
}
