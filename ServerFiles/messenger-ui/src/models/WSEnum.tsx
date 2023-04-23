const WSOutRequest = {
    sendMessage: 'sendMessage',
    login: 'Login',
    sendSeen: 'SendSeen',
    disconnect: 'disconnect',
    getHistory:"GroupHistory",
    createGroup:"CreateGroup"
  };

  const WSInRequest = {
    connect:'connect',
    wsKey:'WSKey',
    login:'Log In',
    readMessage:'ReadMessage',
    history:'History',
    sentMessage:'SentMessage',
    newGroup:'NewGroup',
    newMessage:"NewMessage",
    error:"error",
    disconnect:"disconnect"
  };
  export  {WSOutRequest,WSInRequest};