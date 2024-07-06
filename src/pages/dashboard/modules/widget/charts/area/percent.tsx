import React from "react";
import {
  AreaChart,
  Area,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  ResponsiveContainer,
} from "recharts";

const data = [
  { shift: "Morning", chefs: 40, servers: 60 },
  { shift: "Afternoon", chefs: 50, servers: 50 },
  { shift: "Evening", chefs: 60, servers: 40 },
];

export const ChefServerRatioChart = () => {
  // Calculate percentages for chefs and servers
  const transformedData = data.map((entry) => ({
    shift: entry.shift,
    chefsPercent: (entry.chefs / (entry.chefs + entry.servers)) * 100,
    serversPercent: (entry.servers / (entry.chefs + entry.servers)) * 100,
  }));

  return (
    <div>
      <ResponsiveContainer width='100%' height={400}>
        <AreaChart
          data={transformedData}
          margin={{ top: 20, right: 30, left: 20, bottom: 5 }}
        >
          <CartesianGrid strokeDasharray='3 3' />
          <XAxis dataKey='shift' />
          <YAxis tickFormatter={(value) => `${value}%`} />
          <Tooltip formatter={(value) => `${value}%`} />
          <Area
            type='monotone'
            dataKey='chefsPercent'
            stackId='1'
            stroke='#8884d8'
            fill='#8884d8'
          />
          <Area
            type='monotone'
            dataKey='serversPercent'
            stackId='1'
            stroke='#82ca9d'
            fill='#82ca9d'
          />
        </AreaChart>
      </ResponsiveContainer>
    </div>
  );
};
