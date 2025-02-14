import { useState } from "react";
import "./App.css";
import { StartParsing } from "../wailsjs/go/main/App";

function App() {
  const [filePath, setFilePath] = useState<string>("");
  const [checkInterval, setCheckInterval] = useState<number>(1000);
  const [maximumNumberOfProcessingJobs, setMaximumNumberOfProcessingJobs] =
    useState<number>(1);
  const [maximumExecutionCount, setMaximumExecutionCount] = useState<number>(3);

  const updateFilePath = (e: any) => setFilePath(e.target.value);
  const updateCheckInterval = (e: any) => setCheckInterval(e.target.value);
  const updateMaximumNumberOfProcessingJobs = (e: any) =>
    setMaximumNumberOfProcessingJobs(e.target.value);
  const updateMaximumExecutionCount = (e: any) =>
    setMaximumExecutionCount(e.target.value);

  const [resultText, setResultText] = useState("");
  const updateResultText = (result: string) => setResultText(result);

  const startParsing = () => {
    StartParsing(
      filePath,
      checkInterval,
      maximumNumberOfProcessingJobs,
      maximumExecutionCount,
    ).then(updateResultText);
  };

  return (
    <div id="App" className="input-box">
      Enter directory path:
      <input type="text" value={filePath} onChange={updateFilePath} /> <br />
      Enter check interval:
      <input
        type="number"
        value={checkInterval}
        onChange={updateCheckInterval}
      />
      <br />
      Enter maximum number of processing jobs:
      <input
        type="number"
        value={maximumNumberOfProcessingJobs}
        onChange={updateMaximumNumberOfProcessingJobs}
      />
      <br />
      Enter maximum execution count:
      <input
        type="number"
        value={maximumExecutionCount}
        onChange={updateMaximumExecutionCount}
      />
      <br />
      <button onClick={startParsing}>Start parsing</button>
      <div id="result" className="result">
        {resultText}
      </div>
      {/* <div id="input" className="input-box">
        <input
          id="name"
          className="input"
          onChange={updateName}
          autoComplete="off"
          name="input"
          type="text"
        />
        <button
          className="btn"
          onClick={() => {
            showHelloWorld();
          }}
        >
          Go
        </button>
      </div> */}
    </div>
  );
}

export default App;
