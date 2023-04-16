import Group from "../models/Group";
import { Message } from "../models/Message";
import { User } from "../models/User";
import { BiCheck } from "react-icons/bi";
import { BiCheckDouble } from "react-icons/bi";
import { FcSynchronize } from "react-icons/fc";

function ShowMessage({ message, group, user }: { message: Message, group: Group, user: User }) {
    const dt = new Date(message.sentDate)

    function ProcessState(message:Message)
    {
        switch(message.state){
            case true:
                return <div className="flex w-full justify-end h-0 ml-3"><BiCheckDouble /></div>
            case false:
                return <div className="flex w-full justify-end h-0 ml-3"><BiCheck/></div>
            default:
                return <div className="flex w-full justify-end h-0 ml-3"><FcSynchronize className="animate-spin"/></div>
        }
    }

    return (
        <div className={message.from.number !== user.number || message.from.zone !== user.zone ? "flex w-full justify-start" : "flex w-full justify-end"}>
            <div className="flex flex-col shadow-md w-fit h-fit px-4 py-4 ml-3 mt-3 mr-3 bg-white text-black rounded-md max-w-lg">
                <div>
                    <span className="text-red-600">{group.members.find(member => member.number === message.from.number && member.zone === message.from.zone)?.username}</span>
                </div>
                <div>
                    <span>{message.content}</span>
                </div>

                <div className="flex w-full justify-end">
                        <span className="italic text-[12px] text-gray-400">{dt.toLocaleString().split(", ")[1]}</span>
                </div>
                <div>
                    {message.from.number !== user.number || message.from.zone !== user.zone? "":
                      ProcessState(message)
                    }
                </div>

            </div>
        </div>
    )

}

export default ShowMessage;