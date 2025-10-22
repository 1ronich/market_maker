import { createChart, ColorType, CandlestickSeries } from "lightweight-charts";
import React, { useEffect, useRef, useState } from "react";
import { address, createSolanaRpc } from "@solana/kit";
import { Buffer } from "buffer";
import { binary_to_base58 } from "base58-js";

type RaydiumAMM = {
  BaseVault: string;
  QuoteVault: string;
  BaseMint: string;
  QuoteMint: string;
  LpMint: string;
  OpenOrders: string;
  MarketId: string;
  MarketProgramId: string;
  TargetOrders: string;
  WithdrawQueue: string;
  LpVault: string;
  Owner: string;
};

async function GetTokenPrice(): Promise<number> {
  var b: string = "";
  var q: string = "";

  // Access environment variables
  const rpc_url = process.env.REACT_APP_QUICKNODE_HTTP_SOLANA_DEV;
  console.log("rpc_url:", rpc_url);
  // const rpc_url = "https://api.devnet.solana.com";
  const rpc = createSolanaRpc(rpc_url!);

  /*
  const publicKey = address("97p3acpjRH9kZJ189fRHzVzmWiTSsSoGZeCwFNmHDQhf");
  const accountInfo = await rpc
    .getAccountInfo(publicKey, { encoding: "base64" })
    .send();

  const decoded = await DeserializeRaydiumAMM(accountInfo.value?.data[0]!);

  if (
    decoded.QuoteMint === "So11111111111111111111111111111111111111112" ||
    decoded.QuoteMint === "EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v"
  ) {
    b = //decoded.BaseVault;
    q = decoded.QuoteVault;
  } else {
    b = decoded.QuoteVault;
    q = decoded.BaseVault;
  }*/

  const baseTokenAccount = await rpc
    .getTokenAccountBalance(
      address("8KwXX7228Y6iQw6n9tNWoi5yyY3jm1T6ym4u2ZquvaKT")
    )
    .send();
  const quoteTokenAccount = await rpc
    .getTokenAccountBalance(
      address("4SvfGSq8KxSKHbv6jWNiJoe41R5ih59GAHFrKoRqRFDf")
    )
    .send();

  const quoteSupply =
    Number(quoteTokenAccount.value.amount) /
    Math.pow(10, quoteTokenAccount.value.decimals);
  const baseSupply =
    Number(baseTokenAccount.value.amount) /
    Math.pow(10, baseTokenAccount.value.decimals);

  const price = quoteSupply / baseSupply;

  return price;
}

async function DeserializeRaydiumAMM(b64: string): Promise<RaydiumAMM> {
  const raw = Buffer.from(b64, "base64");

  let offset = 336;

  let BaseVault: Buffer | string = raw.slice(offset, (offset += 32));
  let QuoteVault: Buffer | string = raw.slice(offset, (offset += 32));
  let BaseMint: Buffer | string = raw.slice(offset, (offset += 32));
  let QuoteMint: Buffer | string = raw.slice(offset, (offset += 32));
  let LpMint: Buffer | string = raw.slice(offset, (offset += 32));
  let OpenOrders: Buffer | string = raw.slice(offset, (offset += 32));
  let MarketId: Buffer | string = raw.slice(offset, (offset += 32));
  let MarketProgramId: Buffer | string = raw.slice(offset, (offset += 32));
  let TargetOrders: Buffer | string = raw.slice(offset, (offset += 32));
  let WithdrawQueue: Buffer | string = raw.slice(offset, (offset += 32));
  let LpVault: Buffer | string = raw.slice(offset, (offset += 32));
  let Owner: Buffer | string = raw.slice(offset, (offset += 32));

  BaseVault = binary_to_base58(BaseVault);
  QuoteVault = binary_to_base58(QuoteVault);
  BaseMint = binary_to_base58(BaseMint);
  QuoteMint = binary_to_base58(QuoteMint);
  LpMint = binary_to_base58(LpMint);
  OpenOrders = binary_to_base58(OpenOrders);
  MarketId = binary_to_base58(MarketId);
  MarketProgramId = binary_to_base58(MarketProgramId);
  TargetOrders = binary_to_base58(TargetOrders);
  WithdrawQueue = binary_to_base58(WithdrawQueue);
  LpVault = binary_to_base58(LpVault);
  Owner = binary_to_base58(Owner);

  return {
    BaseVault,
    QuoteVault,
    BaseMint,
    QuoteMint,
    LpMint,
    OpenOrders,
    MarketId,
    MarketProgramId,
    TargetOrders,
    WithdrawQueue,
    LpVault,
    Owner,
  };
}

type SortedData = {
  Highest: number;
  Lowest: number;
};

function SortData(prices_minute: number[]): SortedData {
  var highest: number = 0;
  var lowest: number = 0;

  prices_minute.forEach((p) => {
    if (highest === 0) {
      highest = p;
    }
    if (p > highest) {
      highest = p;
    }

    if (lowest === 0) {
      lowest = p;
    }

    if (p < lowest) {
      lowest = p;
    }
  });

  return { Highest: highest, Lowest: lowest };
}

export const ChartComponent = ({
  dataa,
  colors = {},
}: {
  dataa: any;
  colors?: {
    backgroundColor?: string;
    lineColor?: string;
    textColor?: string;
    areaTopColor?: string;
    areaBottomColor?: string;
  };
}) => {
  const [data, setData] = useState<any[]>([]);
  const {
    backgroundColor = "white",
    lineColor = "#2962FF",
    textColor = "black",
    areaTopColor = "#2962FF",
    areaBottomColor = "rgba(41, 98, 255, 0.28)",
  } = colors;

  const chartContainerRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    /*(async () => {
      const res = await fetch("./data.json")

      const info = await res.json()

      dataa.push(info)

    })()*/
    console.log("data:", data);

    const handleResize = () => {
      chart.applyOptions({
        // width: 400,
        timeScale: {
          timeVisible: true,
          secondsVisible: true,
          barSpacing: 5, // smaller values show more bars and detail
        },
      });
    };
    const chart = createChart(chartContainerRef.current!, {
      layout: {
        background: { type: ColorType.Solid, color: backgroundColor },
        textColor,
      },
      timeScale: {
        timeVisible: true,
        secondsVisible: true,
        barSpacing: 5, // smaller values show more bars and detail
      },
      width: chartContainerRef.current?.clientWidth,
      height: 600,
    });

    chart.timeScale().fitContent();

    const newSeries = chart.addSeries(CandlestickSeries, {
      upColor: "#26a69a",
      downColor: "#ef5350",
      borderVisible: false,
      wickUpColor: "#26a69a",
      wickDownColor: "#ef5350",
    });

    // setData(prevItems => [...prevItems, dataa]);

    console.log("DATA HAHA:", data);
    newSeries.setData(dataa);

    window.addEventListener("resize", handleResize);

    return () => {
      window.removeEventListener("resize", handleResize);

      chart.remove();
    };
  }, [
    dataa,
    backgroundColor,
    lineColor,
    textColor,
    areaTopColor,
    areaBottomColor,
  ]);

  return <div ref={chartContainerRef} />;
};

function GetDate(): string {
  const now = new Date();

  const year = now.getFullYear();
  const month = String(now.getMonth() + 1).padStart(2, "0"); // Months are 0-indexed
  const day = String(now.getDate()).padStart(2, "0");

  const hours = String(now.getHours()).padStart(2, "0");
  const minutes = String(now.getMinutes()).padStart(2, "0");
  const seconds = String(now.getSeconds()).padStart(2, "0");

  // Format: yyyy-mm-dd HH:MM:SS
  const formattedDateTime = `${year}-${month}-${day} ${hours}:${minutes}:${seconds} GMT`;

  return formattedDateTime;
}

const PriceChart = ({ interval }: { interval: number }) => {
  const [prices, setPrices] = useState<any[]>([]);

  var data: number[] = [];
  let i = useRef(1);
  let previous_open = useRef(0);

  useEffect(() => {
    (async () => {
      const res = await fetch("./data.json");
      const info = await res.json();

      setPrices(info);
    })();
  }, []);

  useEffect(() => {
    async function fetchPrice() {
      const price = await GetTokenPrice();
      data.push(price);
    }

    var date: number;
    if (String(i.current).length > 1) {
      date = Number(Math.floor(Date.parse(GetDate()) / 1000));
    } else {
      date = Number(Math.floor(Date.parse(GetDate()) / 1000));
    }

    async function sort() {
      const sorted = SortData(data);
      setPrices([
        ...prices,
        {
          time: date, //new Intl.DateTimeFormat("en-CA").format(new Date()),
          open: previous_open.current, //data[0],
          high: sorted.Highest,
          low: sorted.Lowest,
          close: data[data.length - 1],
        },
      ]);

      i.current++;
      previous_open.current = data[data.length - 1];

      data = [];
    }

    const priceInterval = setInterval(fetchPrice, 100);
    const sortInterval = setInterval(sort, interval);

    return () => {
      clearInterval(priceInterval);
      clearInterval(sortInterval);
    };
  }, [data]);

  return (
    <div>
      <ChartComponent dataa={prices} />
      <button onClick={() => downloadHtml(prices)}>Download json</button>
    </div>
  );
};

const downloadHtml = (data: any) => {
  const jsonData = JSON.stringify(data);
  const blob = new Blob([jsonData], { type: "text/json" });
  const url = URL.createObjectURL(blob);
  const a = document.createElement("a");
  a.href = url;
  a.download = "modified.json";
  a.click();
  URL.revokeObjectURL(url);
};

function MyButton({
  label,
  onClick,
  max,
  min,
  time,
  setMax,
  setMin,
  setTime,
}: {
  label: any;
  onClick: any;
  max: string;
  min: string;
  time: string;
  setMax: React.Dispatch<React.SetStateAction<string>>;
  setMin: React.Dispatch<React.SetStateAction<string>>;
  setTime: React.Dispatch<React.SetStateAction<string>>;
}) {
  return (
    <div
      style={{
        padding: "10px 24px",
        background: "#007bff",
        color: "#fff",
        border: "none",
        borderRadius: "6px",

        display: "flex",
        flexDirection: "row",
        justifyContent: "center", // centers horizontally
        gap: "5px", // space between items
      }}
    >
      <button
        style={{
          cursor: "pointer",
          fontSize: "16px",
        }}
        onClick={onClick}
      >
        {label}
      </button>
      <input
        type="text"
        value={max}
        onChange={(e) => setMax(e.target.value)}
        placeholder="Max price"
        style={{ width: "50%" }} // set width in px, em, %, or ch
      />
      <input
        type="text"
        value={min}
        onChange={(e) => setMin(e.target.value)}
        placeholder="Min price"
        style={{ width: "50%" }} // set width in px, em, %, or ch
      />
      <input
        type="text"
        value={time}
        onChange={(e) => setTime(e.target.value)}
        placeholder="Time range (min)"
        style={{ width: "50%" }} // set width in px, em, %, or ch
      />
    </div>
  );
}

const Maestro = () => {
  const [max, setMax] = useState("");
  const [min, setMin] = useState("");
  const [time, setTime] = useState("");

  useEffect(() => {
    console.log("A CHANGE HAS OCURRED");
  }, [max, min, time]);

  return (
    <div>
      <ManipulationUI
        max={max}
        min={min}
        time={time}
        setMax={setMax}
        setMin={setMin}
        setTime={setTime}
        onClick={() => SendData(max, min, time)}
      ></ManipulationUI>
    </div>
  );
};

const SendData = (max: string, min: string, time: string) => {
  fetch("http://localhost:8080", {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify({ max: max, min: min, time: time }),
  });
};

type ManipulationUIProps = {
  max: string;
  min: string;
  time: string;
  setMax: React.Dispatch<React.SetStateAction<string>>;
  setMin: React.Dispatch<React.SetStateAction<string>>;
  setTime: React.Dispatch<React.SetStateAction<string>>;
  onClick?: () => void; // optional onClick function
};

const ManipulationUI: React.FC<ManipulationUIProps> = ({
  max,
  min,
  time,
  setMax,
  setMin,
  setTime,
  onClick,
}) => {
  return (
    <div>
      <div
        style={{
          display: "flex",
          flexDirection: "row",
          justifyContent: "center", // centers horizontally
          gap: "30px", // space between items
        }}
      >
        <MyButton
          label="send"
          onClick={onClick}
          max={max}
          min={min}
          time={time}
          setMax={setMax}
          setMin={setMin}
          setTime={setTime}
        ></MyButton>
      </div>
    </div>
  );
};

export function App() {
  //  <Maestro></Maestro>             return <ChartComponent {...props, }></ChartComponent>;
  return (
    <div>
      <PriceChart interval={5000} />
    </div>
  );
}
