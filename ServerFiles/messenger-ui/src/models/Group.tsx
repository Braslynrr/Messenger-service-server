import {Message} from "./Message"
import {User} from "./User"

interface Group{
    id:string
	members: User[]
	groupName:string
	description:string
	ischat:boolean
	admins:User[],
	messages?:Message[],
	url?:string
}

export default Group;