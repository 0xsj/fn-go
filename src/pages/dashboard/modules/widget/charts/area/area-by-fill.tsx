import React from "react";
import {
  AreaChart,
  Area,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  ResponsiveContainer,
  Legend,
} from "recharts";

const data = [
  { month: "Jan", employees: 120, tips: 800, satisfaction: 85 },
  { month: "Feb", employees: 125, tips: 850, satisfaction: 87 },
  { month: "Mar", employees: 130, tips: 820, satisfaction: 88 },
  { month: "Apr", employees: 135, tips: 900, satisfaction: 90 },
  { month: "May", employees: 140, tips: 920, satisfaction: 89 },
  { month: "Jun", employees: 145, tips: 950, satisfaction: 91 },
  { month: "Jul", employees: 150, tips: 930, satisfaction: 92 },
];

export const AreaChartFillByValue = () => {
  return (
    <div style={{ width: "100%", height: 400 }}>
      <ResponsiveContainer>
        <AreaChart
          data={data}
          margin={{ top: 20, right: 30, left: 0, bottom: 0 }}
        >
          <CartesianGrid strokeDasharray='3 3' />
          <XAxis dataKey='month' />
          <YAxis />
          <Tooltip />
          <Legend />
          <Area
            type='monotone'
            dataKey='employees'
            stackId='1'
            fill='#8884d8'
            stroke='#8884d8'
          />
          <Area
            type='monotone'
            dataKey='tips'
            stackId='2'
            fill='#82ca9d'
            stroke='#82ca9d'
          />
          <Area
            type='monotone'
            dataKey='satisfaction'
            stackId='3'
            fill='#ffc658'
            stroke='#ffc658'
          />
        </AreaChart>
      </ResponsiveContainer>
    </div>
  );
};
