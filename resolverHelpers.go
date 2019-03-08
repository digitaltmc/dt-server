package main

func makeMeetingItemResolver(agenda []*MeetingItem) []*MeetingItemResolver {
  ret := make([]*MeetingItemResolver, len(agenda))
  for i, v := range agenda {
    ret[i] = &MeetingItemResolver{v}
  }
  return ret
}


