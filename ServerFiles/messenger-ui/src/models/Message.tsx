import {User} from "./User"

interface Message{
	ID:string
	groupID?:string
	from: User
	content:string
	readBy?:{}
	isRead:boolean    
	sentDate:Date
	state?:boolean
}

interface ForwardMessage{
	id?:string
	to:User[]
	from: User
	content:string
}

export type {Message,ForwardMessage};