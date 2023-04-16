import { SetStateAction } from "react";
import Group from "../models/Group";
import {User} from "../models/User";
import { Message } from "../models/Message";


function GroupInfo({ user, group,onSetGroup }: { user: User, group: Group , onSetGroup:(group: SetStateAction<Group|undefined>) => void}){

    group.messages?.sort((x,y)=> new Date(x.sentDate).getTime() - new Date(y.sentDate).getTime())

    function ProcessMessage(msg?:Message):string {
        if(msg!==undefined && msg.content.length>40){
            return `${msg?.content.substring(0,45)}...`
        }

        return msg?.content||""
    }

    function ProcessOwner(msg?:Message):string{
        
        return group.members.find(member=> member.zone=== msg?.from.zone && member.number=== msg?.from.number && member.zone!== user.zone && member.number !== user.number)?.username||"You"
    }

    let friend = group.members.find(member => member.username !== user.username)
    if (group.groupName === "") {
        group.groupName = friend?.username || ""

    }

    if(group.members.length<3)
        group.url=friend?.url

    let newMessages = 0
    for(let message of group.messages||[]){
        if(!message.isRead){
            newMessages++
        }
    }

    return (
        <div className="border border-black flex  border-opacity-50 rounded-md  hover:bg-gray-300" onClick={()=>{onSetGroup(group)}}>
            <img className="h-28 w-28 px-4 py-4 rounded-full " src={group.url} alt="profile-group" />
            <div className="flex flex-col mt-4 grow">
                <span className="font-semibold text-black underline decoration-2 decoration-green-600">{group.groupName}</span>
                <div className="flex">
                    <span className=" text-blue-700 text-[12px]"> {group.members.reduce((members, member) => members + `${member.zone} ${member.number} `, "").replace(`${user.zone} ${user.number}`,group.members.length > 2?"You":"")}</span>
                </div>
                <div>
                   <span className="text-red-600 italic text-sm font-semibold">{group.messages!==undefined && `${group.messages && ProcessOwner(group.messages[group.messages.length-1])}: ` }</span><span className="italic text-sm" >{group.messages &&  ProcessMessage(group.messages[group.messages.length-1])}</span>
                </div>
            </div>
            {newMessages>0 && <div className="bg-red-500 rounded-full h-7 w-7 text-white animate-pulse text-center font-semibold my-1 mx-1">{newMessages}</div>}
        </div>
    )
}

export default GroupInfo;