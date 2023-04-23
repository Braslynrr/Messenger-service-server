import { useRef, useState } from "react";
import { User } from "../models/User";
import { BsGear } from "react-icons/bs";

const UserInfo = ({ user, online ,onUserChanged}: { user?: User, online: boolean, onUserChanged:(user: User)=>void}) => {
    const imageClick = useRef<HTMLInputElement>(null)

    async function UploadFile(files: FileList | null) {
        if (files != null) {
            let header = new Headers()

            header.append("Content-Type", "application/json")

            const buff = await files[0].arrayBuffer()

            const uint8Array = new Uint8Array(buff); // create a new Uint8Array view on the ArrayBuffer

            // convert the Uint8Array to an array of bytes
            const bytes = Array.from(uint8Array);

            let head = {
                method: "POST",
                headers: header,
                body: JSON.stringify({
                    "name": files[0].name,
                    "type": files[0].type,
                    "user": user,
                    "buffer": bytes
                })
            }

            const request = await fetch("/User/ProfileImage", head)

            if (request.status === 200) {
                const data = await request.json()
                onUserChanged({...user!,url:data.url})
            }
        }

    }

    function ClickInputFile() {
        imageClick.current?.click()
    }

    //setUrl(user?.url)
    return (
        <div className="flex overflow-hidden shadow-lg bg-white rounded-lg h-almost-full w-almost-full">
            <div className="flex flex-row mt-4 ml-2 w-full space-x-5" >
                <img alt="profile" className="rounded-full h-32 w-32 py-2 p-2 basis-2/12" onClick={ClickInputFile}
                    src={user?.url}></img>
                <input ref={imageClick} type="file" className="hidden" accept="image/*" onChange={(event) => UploadFile(event.target.files)} />
                <div className="flex flex-col basis-3/6">
                    <div className="underline decoration-2 decoration-sky-500/30 font-black text-md">
                        {user?.username}
                    </div>
                    <div className="text-sky-800 font-semibold">
                        {user?.zone} {user?.number}
                    </div>

                    {online ?
                        <div className="flex space-x-2">
                            <div className="rounded-full bg-green-700 h-3 w-3 mt-2 animate-pulse"></div><span className="underline decoration-2 decoration-green-700">Online</span>
                        </div>
                        :
                        <div className="flex space-x-2">
                            <div className="rounded-full bg-gray-700 h-3 w-3 mt-2 animate-pulse"></div><span className="underline decoration-2 decoration-green-700">Offline</span>
                        </div>
                    }

                    <div className="flex flex-wrap">
                        {user?.state}
                    </div>
                </div>
                <div className="basis-1/12"><BsGear className="h-6 w-6" /></div>
            </div>
        </div>
    )
}

export default UserInfo;