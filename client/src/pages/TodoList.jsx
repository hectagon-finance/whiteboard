import React, { useEffect, useState } from "react";
import elliptic from "elliptic";
import Select from "../common/Select";
import Input from "../common/Input";
import { useForm } from "react-hook-form";
import { yupResolver } from "@hookform/resolvers/yup";
import * as yup from "yup";
import { useLocation } from "react-router-dom";
import generateRandomNumber from "../function/generateRandomNumber";
import * as WebSocket from "websocket";
import { Buffer } from "buffer";

const EC = elliptic.ec;
const ec = new EC("p256");

const Command = ["Create", "Start", "Stop", "Pause", "Finish", "Assign"];

const TodoList = () => {
  const schema = yup.object().shape({
    id: yup.string().required("Field is invalid"),
    description: yup.string().required("Field is invalid"),
    title: yup.string().required("Field is invalid"),
  });
  const {
    register,
    handleSubmit,
    formState: { errors },
  } = useForm({
    mode: "onChange",
    resolver: yupResolver(schema),
  });
  const location = useLocation();
  const [publicKey, setPublicKey] = useState(null);
  const [walletAddress, setWalletAddress] = useState(null);
  const [selectedCommand, setSelectedCommand] = useState("Create");
  const [socket, setSocket] = useState(null);
  const [isConnected, setIsConnected] = useState(false);

  useEffect(() => {
    const ws = new WebSocket.w3cwebsocket("ws://localhost:9000/ws");
    setSocket(ws);
    if (location.state) {
      setPublicKey(location.state.publicKey);
      setWalletAddress(location.state.walletAddress);
    }

    ws.onopen = () => {
      setIsConnected(true);
      console.log("WebSocket connected");
    };

    ws.onclose = () => {
      setIsConnected(false);
      console.log("WebSocket disconnected");
    };

    return () => {
      ws.close();
    };
  }, [location]);

  console.log("privateKey: " + location.state.privateKey);
  console.log("publicKey: " + publicKey);
  console.log("walletAddress: " + walletAddress);

  const handleSubmitForm = (valueFields) => {
    const keyPair = ec.keyFromPrivate(location.state.privateKey, "hex");

    let instructionData = {
      Id: valueFields.id,
      description: valueFields.description,
      title: valueFields.title,
      From: walletAddress.toString(),
    };
    const instructionDataJsonStr = JSON.stringify(instructionData);
    const instructionDataBase64Str = Buffer.from(
      instructionDataJsonStr
    ).toString("base64");

    const data = JSON.stringify({
      C: selectedCommand,
      Data: instructionDataBase64Str,
    });

    const encoder = new TextEncoder();
    const byteArray = encoder.encode(data);
    const signature = keyPair.sign(byteArray);

    const signatureStr =
      signature.s.toString(16).padStart(64, "0") +
      signature.r.toString(16).padStart(64, "0");

    let message = JSON.stringify({
      type: "transaction",
      from: "client",
      transactionId: generateRandomNumber(),
      publicKey: publicKey,
      signature: signatureStr,
      data: data,
    });

    console.log("message: " + message);

    if (isConnected && socket) {
      socket.send(message);
    } else {
      console.error("WebSocket not connected");
    }
  };

  const handleCommandChange = (newCommand) => {
    setSelectedCommand(newCommand);
    console.log("Selected command in TodoList: " + newCommand);
  };

  return (
    <div className="flex flex-col items-center justify-center h-screen">
      <Select
        listCommands={Command}
        selectedCommand={selectedCommand}
        onCommandChange={handleCommandChange}
      />
      <div className="bg-white p-10 rounded-lg shadow md:w-3/4 mx-auto lg:w-1/2 mt-[30px]">
        <form onSubmit={handleSubmit(handleSubmitForm)}>
          {selectedCommand === "Create" && (
            <div>
              <Input
                label={"Id"}
                id="id"
                register={register("id")}
                message={errors?.id?.message}
              />
              <Input
                label={"Description"}
                id="description"
                register={register("description")}
                message={errors?.description?.message}
              />
              <Input
                label={"Title"}
                id="title"
                register={register("title")}
                message={errors?.title?.message}
              />
            </div>
          )}
          {/* <Input
            label={"Repo name"}
            id="repo"
            register={register("repoName")}
            message={errors?.repoName?.message}
          /> */}
          <button className="block w-full bg-blue-500 text-white font-bold p-4 rounded-lg">
            Submit
          </button>
        </form>
      </div>
    </div>
  );
};

export default TodoList;
