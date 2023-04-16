import { IconContext } from "react-icons/lib";
import Group from "../models/Group";
import {User,equalUser} from "../models/User";
import ShowMessage from "./ShowMessage";
import { IoIosSend } from "react-icons/io";
import { SetStateAction, useState } from "react";
import { Socket } from "socket.io-client";
import { DefaultEventsMap } from "socket.io/dist/typed-events";
import { encryptObject,Key } from "../Utils/Utils";
import { ForwardMessage, Message} from "../models/Message";
import { WSOutRequest } from "../models/WSEnum";

function Chat({user,group,client,onSendMessage,onMessageChanged}:{user:User,group:Group,client:Socket<DefaultEventsMap, DefaultEventsMap>, onSendMessage:(group: SetStateAction<Group|undefined>) => void, onMessageChanged:(message: Message)=>void}) {
    const [content,setContent] = useState<string>("")

    group.messages?.sort((x,y)=> new Date(x.sentDate).getTime() - new Date(y.sentDate).getTime())

    function sendSeens(){

        if(group.messages!==undefined){
            for(let message of group.messages){

                if(!message.isRead){
                    client.emit(WSOutRequest.sendSeen, message.ID)
                    message.isRead = true
                    onMessageChanged(message)
                }
            }
    
        }

    }

    function sendMessage(){
        setContent("")
        let newMessage:ForwardMessage = {id:group.id ,from:user,content:content,to:group.members.filter(member=> !equalUser(user,member))}
        let encrypted = encryptObject(newMessage,Key)

        if(group.messages===undefined)
            group.messages=[]

        group.messages.push({ID:"",from:user,content:content,isRead:true,sentDate:new Date()})

        let newgroup = {...group}
        onSendMessage(newgroup)

        client.emit(WSOutRequest.sendMessage, encrypted)
    }

    return (
        <IconContext.Provider value={{color:"#00ffff"}}>
        <div className="flex flex-col items-center justify-center h-full w-full shadow-lg bg-gray-900 text-white" onClick={sendSeens}>
            <div className="basis-1/12 bg-white w-full text-black hover:bg-gray-300">
                <div className="flex w-full">
                    <div className="flex flex-col m-auto text-center">
                        <span className=" font-semibold text-2xl">{group.groupName}</span>
                        <span className="text-blue-600">{group.members.reduce((members, member) => members + `${member.zone} ${member.number} `, "").replace(`${user.zone} ${user.number}`,group.members.length > 2?"You":"")}</span>
                    </div>
                </div>
            </div>
            <div className="flex flex-col basis-10/12 w-full overflow-y-auto">
                {group.messages && group.messages.map(message => <ShowMessage user={user} group={group} message={message}/>)}
            </div>
            <div className="basis-1/12">
                <div className="flex space-x-3">
                    <div>
                        <textarea value={content} onChange={(event)=> setContent(event.target.value)}
                        className="border px-8 bg-transparent rounded-2xl border-white max-h-14 min-w-full shadow appearance-none border-opacity-50 focus:border-blue-500 focus:outline-none focus:shadow-outline"></textarea>  
                    </div>
                    <div><IoIosSend className="h-10 w-10 mt-2 rounded-full hover:bg-transparent hover:animate-pulse" onClick={sendMessage}/></div>
                </div>
            </div>
        </div>
    </IconContext.Provider>
    )
}

export default Chat;