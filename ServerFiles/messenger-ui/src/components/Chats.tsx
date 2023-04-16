import { Fragment, SetStateAction, useCallback, useState } from "react";
import Group from "../models/Group";
import {User} from "../models/User";
import GroupInfo from "./GroupInfo";
import NewMessageModal from "./NewMessageModal";

function Chats({user,groups,onSetGroup}:{user:User,groups:Group[],onSetGroup:(group: SetStateAction<Group|undefined>) => void}){

const [showModal, setShowModal] = useState<boolean>(false);

const updateModal = useCallback((group: SetStateAction<boolean>) => {
    setShowModal(group)
}, [])

return(
    <Fragment>
    {showModal && <NewMessageModal onSetGroup={onSetGroup}  user={user}  setOpenModal={updateModal} />}
    <div className="h-almost-full w-almost-full shadow-lg bg-white rounded-lg z-50">
        <div className="h-full w-full flex flex-col">
            <div className="flex">
               <button className="flex bg-blue-500 text-white rounded-md mt-2 mb-2 mx-auto px-2 hover:bg-blue-400" onClick={()=> setShowModal(!showModal)}>Start Chat</button> 
               <input key="find" className="mx-auto rounded-md border border-blue-500 mt-3 mb-3" placeholder=" Search"></input>
            </div>
            <div className="grow overscroll-y-auto">{}
                {groups.map(group => <GroupInfo key={group.id} user={user} group={group} onSetGroup={onSetGroup} />)}
            </div>
        </div>
    </div>

    </Fragment>
)};

export default Chats;