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
  const schemaCreate = yup.object().shape({
    idCreate: yup.string().required("Field is invalid"),
    description: yup.string().required("Field is invalid"),
    title: yup.string().required("Field is invalid"),
  });

  const schemaStart = yup.object().shape({
    idStart: yup.string().required("Field is invalid"),
    estdaytofinish: yup.string().required("Field is invalid"),
  });

  const schemaStop = yup.object().shape({
    idStop: yup.string().required("Field is invalid"),
    reason: yup.string().required("Field is invalid"),
  });

  const schemaPause = yup.object().shape({
    idPause: yup.string().required("Field is invalid"),
    estwaitday: yup.string().required("Field is invalid"),
  });
  const schemaFinish = yup.object().shape({
    idFinish: yup.string().required("Field is invalid"),
    congratmessage: yup.string().required("Field is invalid"),
  });

  const schemaAssign = yup.object().shape({
    idAssign: yup.string().required("Field is invalid"),
    assignTo: yup.string().required("Field is invalid"),
  });
  const [schema, setSchema] = useState(schemaCreate);
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
  const [base64Data, setBase64Data] = useState(null);

  useEffect(() => {
    switch (schema) {
      case "Create":
        setSchema(schemaCreate);
        break;
      case "Start":
        setSchema(schemaStart);
        break;
      case "Stop":
        setSchema(schemaStop);
        break;
      case "Pause":
        setSchema(schemaPause);
        break;
      case "Finish":
        setSchema(schemaFinish);
        break;
      case "Assign":
        setSchema(schemaAssign);
        break;
      default:
        break;
    }
  }, [schema]);

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

    fetchData();

    ws.onclose = () => {
      setIsConnected(false);
      console.log("WebSocket disconnected");
    };

    return () => {
      ws.close();
    };
  }, [location, selectedCommand]);

  const fetchData = async () => {
    try {
      const response = await fetch("http://localhost:18080/get");
      const data = await response.text();
      console.log("data", data);
      setBase64Data(data);
    } catch (error) {
      console.error("Error fetching data:", error);
    }
  };

  const handleSubmitForm = (valueFields) => {
    const ws = new WebSocket.w3cwebsocket("ws://localhost:9000/ws");
    setSocket(ws);

    console.log("===============send");
    const keyPair = ec.keyFromPrivate(location.state.privateKey, "hex");

    let instructionData;

    switch (selectedCommand) {
      case "Create":
        instructionData = {
          Id: valueFields.idCreate,
          Desc: valueFields.description,
          Title: valueFields.title,
          From: walletAddress.toString(),
        };
        console.log(instructionData);
        break;
      case "Start":
        console.log("=============Start");
        instructionData = {
          Id: valueFields.idStart,
          EstDayToFinish: parseInt(valueFields.estdaytofinish),
          From: walletAddress.toString(),
        };
        console.log(instructionData);
        break;
      case "Stop":
        instructionData = {
          Id: valueFields.idStop,
          Reason: valueFields.reason,
          From: walletAddress.toString(),
        };
        console.log(instructionData);
        break;
      case "Pause":
        instructionData = {
          Id: valueFields.idPause,
          EstWaitDay: parseInt(valueFields.estwaitday),
          From: walletAddress.toString(),
        };
        console.log(instructionData);
        break;
      case "Finish":
        instructionData = {
          Id: valueFields.idFinish,
          CongratMessage: valueFields.congratmessage,
          From: walletAddress.toString(),
        };
        console.log(instructionData);
        break;
      case "Assign":
        instructionData = {
          Id: valueFields.idAssign,
          AssignTo: valueFields.assignTo,
          From: walletAddress.toString(),
        };
        console.log("instructionData", instructionData);
        break;
      default:
        break;
    }

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

    // Delay fetchData() by 2 seconds
    setTimeout(() => {
      fetchData();
    }, 2000);
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
                register={register("idCreate")}
                message={errors?.idCreate?.message}
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
          {selectedCommand === "Start" && (
            <div>
              <Input
                label={"Id"}
                register={register("idStart")}
                message={errors?.idStart?.message}
              />
              <Input
                label={"Est Day To Finish"}
                id="estdaytofinish"
                register={register("estdaytofinish")}
                message={errors?.estdaytofinish?.message}
              />
            </div>
          )}
          {selectedCommand === "Stop" && (
            <div>
              <Input
                label={"Id"}
                register={register("idStop")}
                message={errors?.idStop?.message}
              />
              <Input
                label={"Reason"}
                id="reason"
                register={register("reason")}
                message={errors?.reason?.message}
              />
            </div>
          )}
          {selectedCommand === "Pause" && (
            <div>
              <Input
                label={"Id"}
                register={register("idPause")}
                message={errors?.idPause?.message}
              />
              <Input
                label={"Est Wait Day"}
                id="estwaitday"
                register={register("estwaitday")}
                message={errors?.estwaitday?.message}
              />
            </div>
          )}
          {selectedCommand === "Finish" && (
            <div>
              <Input
                label={"Id"}
                register={register("idFinish")}
                message={errors?.idFinish?.message}
              />
              <Input
                label={"Congrat Message"}
                id="congratmessage"
                register={register("congratmessage")}
                message={errors?.congratmessage?.message}
              />
            </div>
          )}
          {selectedCommand === "Assign" && (
            <div>
              <Input
                label={"Id"}
                register={register("idAssign")}
                message={errors?.idAssign?.message}
              />
              <Input
                label={"Assign"}
                id="assign"
                register={register("assignTo")}
                message={errors?.assignTo?.message}
              />
            </div>
          )}
          <button className="block w-full bg-blue-500 text-white font-bold p-4 rounded-lg">
            Submit
          </button>
        </form>
      </div>
      <div>{base64Data && <p>{base64Data}</p>}</div>
    </div>
  );
};

export default TodoList;
