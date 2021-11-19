import runExperiment from './experiment';

const looseJSONParse = (str) => {
    try {
        return JSON.parse(str);
    } catch (error) {
        console.error(error);
    }
};

const init = () => {
    console.log("[revcor] version 0.2");

    const wsProtocol = window.location.protocol === "https:" ? "wss" : "ws";
    const pathPrefixhMatch = /(.*)xp/.exec(window.location.pathname);
    // depending on APP_WEB_PREFIX, signaling endpoint may be located at /ws or /prefix/ws
    const pathPrefix = pathPrefixhMatch[1];
    const signalingUrl = `${wsProtocol}://${window.location.host}${pathPrefix}ws`;
    const ws = new WebSocket(signalingUrl);

    ws.onopen = () => {
        const { experimentId, participantId } = window.state;
        ws.send(            
            JSON.stringify({
                kind: "join",
                payload: JSON.stringify({ experimentId, participantId }),
            })
        );
    };

    ws.onmessage = async (event) => {
        let message = looseJSONParse(event.data);
        
        if (message.kind === "init") {
            message.payload.wording = looseJSONParse(message.payload.wording);
            runExperiment(message.payload, ws);
        }
    };

}

document.addEventListener('DOMContentLoaded', init);
