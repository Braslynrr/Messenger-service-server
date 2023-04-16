import { useState, useEffect, useCallback, SetStateAction } from 'react';
import { useNavigate } from 'react-router-dom';
import { User } from '../models/User';
import { io } from "socket.io-client";
import UserInfo from '../components/UserInfo';
import Chats from '../components/Chats';
import Chat from '../components/Chat';
import Group from '../models/Group';
import { Message } from '../models/Message';
import { decryptText, setKey } from '../Utils/Utils';
import CryptoJS from 'crypto-js';
import { WSInRequest, WSOutRequest } from '../models/WSEnum';


const socketUrl = 'http://localhost:8080/';
const client = io(socketUrl, {
    autoConnect: false,
    transports: ['websocket'],
})

function Messenger() {

    const navigate = useNavigate()
    const [key, setWSKey] = useState<CryptoJS.lib.WordArray>(null!)
    const [user, setUser] = useState<User>(new User())
    const [groups, setGroups] = useState<Group[]>([])
    const [selectedGroup, setSelectedGroup] = useState<Group>()
    const [online, setOnline] = useState<boolean>(false)
    const HandleReadMessage = useCallback((message: Message) => {
        function onMessageChange(incomingMsg: Message) {
            let newGroups = [...groups]

            let tempGroup = newGroups.find(group => group.id === incomingMsg.groupID)

            if (tempGroup !== undefined) {

                if (tempGroup.messages !== undefined) {

                    
                    let index = tempGroup.messages.findIndex(msg => msg.ID === incomingMsg.ID)

                    if (index !== undefined) {
                        incomingMsg.isRead = true

                        tempGroup.messages![index] = incomingMsg

                        setGroups(newGroups)
                    }
                }


            }
        }
        onMessageChange(message)
    }, [groups])
    const HandleSetGroup = useCallback((group: SetStateAction<Group | undefined>) => setSelectedGroup(group), [])



    function Connected() {
        setOnline(true)
    }


    function Login(user: User) {
        if (user.url === "" || user.url === undefined)
            user.url = "https://external-content.duckduckgo.com/iu/?u=https%3A%2F%2Fwww.w3schools.com%2Fw3css%2Fimg_avatar3.png&f=1&nofb=1&ipt=8be56eab19d1bb80223b87b32641854d9d88cb4ddf7e466a01be19d39332c65e&ipo=images"

        setUser(user)
    }

    function AddKey(key: string) {
        let wskey: CryptoJS.lib.WordArray = CryptoJS.enc.Hex.parse(key)
        setWSKey(wskey)
        setKey(wskey)
    }

    function error(err: any) {
        console.log(err.Error)
    }



    useEffect(
        () => {

            function NewGroup(group: Group) {
                let Newgroups: Group[] = [group, ...groups]
                setGroups(Newgroups)
            }

            function SentMessage(SentMessage: any) {
                let newGroups = [...groups]
                let message = newGroups.find(group => group.id === SentMessage.message.groupid)?.messages?.find(msg => msg.ID === "")
                if (message !== undefined) {
                    message!.ID = SentMessage.message.id
                    message.state = false
                }

                setGroups(newGroups)
            }

            function readMessage(encryptedMessage: string) {

                let message = decryptText<Message>(encryptedMessage, key)

                let newGroups = [...groups]

                let tempGroup = newGroups.find(group => group.id === message.groupID)

                if (tempGroup !== undefined) {

                    let index = tempGroup!.messages?.findIndex(msg => msg.ID === message.ID)!

                    if (message.readBy !== undefined) {
                        if (Object.keys(message.readBy).length === tempGroup.members.length - 1) {
                            message.state = true
                        } else {
                            message.state = false
                        }
                    }

                    tempGroup!.messages![index] = message

                    setGroups(newGroups)
                }
            }


            function NewMessage(encryptedMessage: string) {
                let message = decryptText<Message>(encryptedMessage, key)

                let newGroups = [...groups]

                let tempGroup = newGroups.find(group => group.id === message.groupID)

                if (groups !== undefined) {
                    if (tempGroup?.messages === undefined) {

                        tempGroup!.messages = [message]

                    } else {

                        tempGroup!.messages.push(message)
                    }
                }

                setGroups(newGroups)
            }

            client.on(WSInRequest.newGroup, NewGroup)
            client.on(WSInRequest.sentMessage, SentMessage)

            if (key !== null) {

                client.on(WSInRequest.readMessage, readMessage)
                client.on(WSInRequest.newMessage, NewMessage)
            }

            return () => {
                client.off(WSInRequest.newGroup, NewGroup)
                client.off(WSInRequest.sentMessage, SentMessage)
                if (key !== null) {
                    client.off(WSInRequest.readMessage, readMessage)
                    client.off(WSInRequest.newMessage, NewMessage)
                }
            }


        }, [groups, key])

    useEffect(() => {

        async function ManageGroups() {

            let header = new Headers()

            header.append("Content-Type", "application/json")

            let req = {
                method: "GET",
                headers: header,
            }

            let request = await fetch(`/Groups/${user?.zone}/${user?.number}/`, req)

            let groups: Group[] = await request.json()

            if (groups !== undefined && groups !== null) {

                header = new Headers()

                header.append("Content-Type", "text/plain")


                for (let group of groups) {
                    let body = {
                        method: "POST",
                        headers: header,
                        body: JSON.stringify({ socketID: client.id, time: new Date(), ID: group.id })
                    }
                    let req = await fetch("/Groups/Messages", body);

                    let data = await req.text()

                    let messages: Message[] = []

                    if (data !== "")
                        messages = decryptText<Message[]>(data, key)

                    messages.forEach(msg => {
                        if (msg.readBy !== undefined) {

                            if (Object.keys(msg.readBy).length === group.members.length - 1) {
                                msg.state = true
                            } else {
                                msg.state = false
                            }
                        } else {
                            msg.state = false
                        }
                    })

                    group.messages = messages


                }
                setGroups(groups)

            }

        }

        function Logout(error: any) {
            client.disconnect()
            navigate("/LogIn")
        }

        if (sessionStorage.token == null)
            navigate("/")
        client.on(WSInRequest.connect, Connected)
        client.on(WSInRequest.login, Login)
        client.on(WSInRequest.wsKey, AddKey)
        client.on("errorLogin", Logout)
        client.on(WSInRequest.error, error)



        if (!client.connected)
            client.connect()

        if (user.zone === "")
            client.emit(WSOutRequest.login, sessionStorage.token)

        if (key !== null) {

            if (user.zone !== "")
                ManageGroups()
        }

        return () => {
            client.off(WSInRequest.connect, Connected)
            client.off(WSInRequest.login, Login)
            client.off(WSInRequest.wsKey, AddKey)
            client.off("errorLogin", Logout)
            client.off(WSInRequest.error, error)

        };


    }, [key, navigate, user])

    return (<div className="h-screen w-screen text-black flex flex-row">
        <div className=" basis-full md:basis-1/4 flex flex-col space-y-1">
            <div className='justify-center items-center h-full basis-2/12 ml-1'>
                <UserInfo user={user} online={online} />
            </div>
            <div className='grow justify-center items-center ml-1'>
                <Chats user={user} groups={groups} onSetGroup={HandleSetGroup} />
            </div>
        </div>
        <div className="md:grow justify-center items-center">
            {key &&
                selectedGroup ? <Chat client={client} user={user} group={selectedGroup} onSendMessage={HandleSetGroup} onMessageChanged={HandleReadMessage} /> : <div className='flex bg-gray-900 w-full h-full text-white text-center justify-center'><div className='m-auto text-2xl'>Welcome to Messenger Service</div></div>
            }
        </div>

    </div>)
}

export default Messenger;