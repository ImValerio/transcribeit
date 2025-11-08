"use client";
import Image from "next/image";
import { useState } from "react";

export default function Home() {
  const hostname = "http://localhost:8080";
  const [fileName, setFileName] = useState("");
  const [transcript, setTranscript] = useState("");
  const [isLoading, setIsLoading] = useState(false);
  const [isUploaded, setIsUploaded] = useState(false);

  const handleSubmit = async (event: React.FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    const formData = new FormData(event.currentTarget);
    const file = formData.get("file") as File;
    if (!file || file.name === "") {
      alert("Please select a file");
      return;
    }
    setIsUploaded(false)
    setTranscript("")

    try {
      setIsLoading(true);
      const response = await fetch(`${hostname}/transcribe`, {
        method: "POST",
        body: formData,
      });

      const { audio_id } = await response.json();
      setIsUploaded(true)

      await pollTranscription(audio_id);

    } catch (error) {
      console.error(error);
    }
  };


  async function pollTranscription(id: string) {
    const res = await fetch(`${hostname}/transcribe/${id}`);
    const data = await res.json();

    if (data.status === "pending") {
      console.log("Still processing...");
      setTimeout(() => pollTranscription(id), 2000); // try again in 2s
    } else {
      setIsLoading(false)
      setTranscript(data.transcription)
    }
  }

  return (
    <div className="font-sans flex flex-col justify-between items-center justify-items-center h-full p-8 pb-20 gap-16 sm:p-20">
      {/* <h1 className="text-5xl md:text-7xl tracking-widest animate__animated animate__fadeInDown">
        Transcribe
        <span className="font-bold ">IT</span>
      </h1> */}
      <img width={300} height={150} alt={"logo"} src={"transcribeit_banner.png"} className="animate__animated animate__fadeInDown scale-150 md:scale-130" />
      <form
        className="flex flex-col gap-[16px] border-slate-300 border-2 p-4"
        id="upload-audio-form"
        onSubmit={handleSubmit}
        onChange={(event) => {
          const file = event.target.files[0];
          if (file) {
            setFileName(file.name);
          }
        }}
      >
        <input
          type="file"
          id="file"
          accept=".mp3,.m4a,.wav"
          name="file"
          className="hidden"
        />
        <label
          htmlFor="file"
          className="text-2xl font-medium text-gray-500 flex justify-center items-center"
        >
          {!fileName && (
            <Image
              aria-hidden
              src="/file-add.svg"
              alt="File icon"
              width={25}
              height={25}
              className="invert brightness-0"
            />
          )}
          Select audio file
        </label>
        {fileName && (
          <p className="text-sm text-gray-500 flex justify-center items-center">
            <Image
              aria-hidden
              src="/file-basic.svg"
              color="white"
              alt="File icon"
              width={16}
              height={16}
              className="invert brightness-0"
            />
            {fileName}
          </p>
        )}
        <input
          className="transcribe-btn text-white font-bold py-2 px-4 rounded cursor-pointer transition-all duration-300 ease-in-out"
          type="submit"
          value="Transcribe"
        />
      </form>

      <main className="flex flex-col gap-[32px] row-start-2 items-center justify-center animate__animated animate__fadeIn">
        {isLoading && (

          <div className="flex  justify-center items-center">
            <ul className="p-3">
              <li className="flex justify-center items-center">Upload {isUploaded ? "✅" :
                <div className="loader animate__animated animate__zoomIn w-max-3"></div>
              } </li>
              <li className="flex justify-center items-center">
                Transcription {transcript ? "✅" :
                  <div className="loader animate__animated animate__zoomIn w-max-3 mx-5"></div>
                } </li>
            </ul>
          </div>
        )}
        {transcript && !isLoading && (
          <div className="flex flex-col justify-center items-center w-full md:w-max-xl animate__animated animate__fadeInUp">
            <header className="flex w-full justify-end bg-primary p-2">
              <p
                className="cursor-pointer"
                onClick={() => navigator.clipboard.writeText(transcript)}
              >
                Copy to clipboard
              </p>
            </header>
            <textarea
              cols={70}
              className="text-sm text-gray-500 flex justify-center items-center w-full md:w-max-xl max-h-49 overflow-y-auto word-wrap break-normal scroll-auto m-0 "
              readOnly={true}
              value={transcript}
            ></textarea>
          </div>
        )}
      </main>

      <footer className="row-start-4 flex gap-[24px] flex-wrap items-center justify-center">
        <a
          className="flex items-center gap-2 hover:underline hover:underline-offset-4"
          href="https://github.com/imvalerio/transcribeit"
          target="_blank"
          rel="noopener noreferrer"
        >
          <Image
            aria-hidden
            src="/file.svg"
            alt="File icon"
            width={16}
            height={16}
          />
          Source code
        </a>
        <a
          className="flex items-center gap-2 hover:underline hover:underline-offset-4"
          href="https://valeriovalletta.it"
          target="_blank"
          rel="noopener noreferrer"
        >
          <Image
            aria-hidden
            src="/globe.svg"
            alt="Globe icon"
            width={16}
            height={16}
          />
          About me
        </a>
      </footer>
    </div>
  );
}
